// Package main 是 Gateway 服务的入口，初始化配置、发现、路由与 WebSocket 服务并启动 HTTP 服务器。
package main

import (
	"IM/pkg/auth"
	"IM/pkg/config"
	"IM/pkg/discovery"
	"IM/pkg/logger"
	"IM/services/gateway/handler"
	"IM/services/gateway/middleware"
	"IM/services/gateway/proxy"
	"IM/services/gateway/ws"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.LoadGatewayConfig(getConfigPath())
	if err != nil {
		fmt.Printf("Using default config: %v\n", err)
		cfg = config.DefaultGatewayConfig()
	}

	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	resolver, err := discovery.NewResolver(cfg.Etcd.Endpoints, cfg.Etcd.DialTimeout)
	if err != nil {
		logger.Fatalw("Failed to create resolver", "component", "gateway_cmd", "err", err)
	}
	defer resolver.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	jwtUtil := auth.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.Expire)

	serviceProxy := proxy.NewServiceProxy(resolver, cfg)
	hub := ws.NewHub()
	go hub.Run()

	wsHandler := ws.NewHandler(hub, jwtUtil, redisClient)

	router := gin.Default()

	router.Use(middleware.CORS())
	router.Use(middleware.Logging())
	router.Use(middleware.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "gateway",
		})
	})

	api := router.Group("/api/v1")
	{
		api.POST("/auth/register", handler.Register(serviceProxy))
		api.POST("/auth/login", handler.Login(serviceProxy, jwtUtil))

		protected := api.Group("")
		protected.Use(middleware.Auth(jwtUtil))
		{
			protected.GET("/users/:id", handler.GetUser(serviceProxy))
			protected.GET("/users/friends", handler.GetFriends(serviceProxy))
			protected.POST("/users/friends", handler.AddFriend(serviceProxy))
			protected.POST("/users/friend_requests/accept", handler.AcceptFriendRequest(serviceProxy))
			protected.GET("/users/friend_requests", handler.GetPendingFriendRequests(serviceProxy))

			protected.POST("/groups", handler.CreateGroup(serviceProxy))
			protected.GET("/groups/:id", handler.GetGroup(serviceProxy))
			protected.GET("/groups/:id/members", handler.GetGroupMembers(serviceProxy))
			protected.GET("/users/groups", handler.GetUserGroups(serviceProxy))
			protected.POST("/groups/join", handler.JoinGroup(serviceProxy))
			protected.GET("/groups/join/requests", handler.GetPendingGroupJoinRequests(serviceProxy))
			protected.GET("/groups/join/accept", handler.AcceptGroupJoinRequest(serviceProxy))
			protected.DELETE("/groups/:id/leave", handler.LeaveGroup(serviceProxy))

			protected.POST("/messages", handler.SendMessage(serviceProxy))
			protected.GET("/messages/offline", handler.GetOfflineMessages(serviceProxy))
			protected.PUT("/messages/:id/read", handler.MarkAsRead(serviceProxy))
			protected.PUT("/messages/read/all", handler.MarkAllAsRead(serviceProxy))
			protected.GET("/messages/unread/count", handler.GetUnreadCount(serviceProxy))
		}
	}

	router.GET("/ws", func(c *gin.Context) {
		wsHandler.HandleWebSocket(c.Writer, c.Request)
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Server.HTTPPort,
		Handler: router,
	}

	go func() {
		logger.Infow("API Gateway starting", "component", "gateway_cmd", "http_port", cfg.Server.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalw("Failed to start server", "component", "gateway_cmd", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infow("Shutting down server...", "component", "gateway_cmd")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorw("Server shutdown error", "component", "gateway_cmd", "err", err)
	}

	logger.Infow("Server stopped", "component", "gateway_cmd")
}

// getConfigPath 返回配置文件路径，优先使用环境变量 CONFIG_PATH，否则返回默认路径。
func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "configs/gateway.yaml"
}
