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
	// 加载配置
	cfg, err := config.LoadGatewayConfig(getConfigPath())
	if err != nil {
		fmt.Printf("Using default config: %v\n", err)
		cfg = config.DefaultGatewayConfig()
	}
	// 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	// 初始化服务发现
	resolver, err := discovery.NewResolver(cfg.Etcd.Endpoints, cfg.Etcd.DialTimeout)
	if err != nil {
		logger.Fatalw("Failed to create resolver", "component", "gateway_cmd", "err", err)
	}
	defer resolver.Close()
	// 初始化 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()
	// 初始化 JWT 工具
	jwtUtil := auth.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.Expire)
	// 初始化服务代理
	serviceProxy := proxy.NewServiceProxy(resolver, cfg)

	// 初始化 WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()
	defer hub.Stop()

	// 初始化 WebSocket 处理器
	wsHandler := ws.NewHandler(hub, jwtUtil, redisClient)

	// 初始化消息消费者（从 RabbitMQ 拉取消息推送给在线用户）
	var msgConsumer *ws.MessageConsumer
	if cfg.RabbitMQ.URL != "" {
		queueName := cfg.RabbitMQ.QueuePrefix + "gateway"
		msgConsumer, err = ws.NewMessageConsumer(hub, cfg.RabbitMQ.URL, cfg.RabbitMQ.Exchange, queueName)
		if err != nil {
			logger.Warnw("Failed to create message consumer, messages will not be pushed via WebSocket", "component", "gateway_cmd", "err", err)
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if err := msgConsumer.Start(ctx); err != nil {
				logger.Warnw("Failed to start message consumer", "component", "gateway_cmd", "err", err)
			} else {
				logger.Infow("Message consumer started", "component", "gateway_cmd", "queue", queueName)
			}
			defer msgConsumer.Close()
		}
	}

	// 初始化路由
	router := gin.Default()
	// 初始化中间件
	router.Use(middleware.CORS())
	router.Use(middleware.Logging())
	router.Use(middleware.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":       "healthy",
			"service":      "gateway",
			"online_users": hub.OnlineCount(),
		})
	})
	// API 路由表
	api := router.Group("/api/v1")
	{
		api.POST("/auth/register", handler.Register(serviceProxy))
		api.POST("/auth/login", handler.Login(serviceProxy, jwtUtil))
		// 受保护路由
		protected := api.Group("")
		protected.Use(middleware.Auth(jwtUtil))
		{
			protected.GET("/users/:id", handler.GetUser(serviceProxy))
			protected.GET("/users/friends", handler.GetFriends(serviceProxy))
			protected.POST("/users/friends", handler.AddFriend(serviceProxy))
			protected.DELETE("/users/friends/:friend_id", handler.RemoveFriend(serviceProxy))
			protected.POST("/users/friend_requests/accept", handler.AcceptFriendRequest(serviceProxy))
			protected.GET("/users/friend_requests", handler.GetPendingFriendRequests(serviceProxy))

			protected.POST("/groups", handler.CreateGroup(serviceProxy))
			protected.GET("/groups/:id", handler.GetGroup(serviceProxy))
			protected.GET("/groups/:id/members", handler.GetGroupMembers(serviceProxy))
			protected.DELETE("/groups/:id/members/:member_id", handler.RemoveGroupMember(serviceProxy))
			protected.GET("/users/groups", handler.GetUserGroups(serviceProxy))
			protected.POST("/groups/join", handler.JoinGroup(serviceProxy))
			protected.GET("/groups/join/requests", handler.GetPendingGroupJoinRequests(serviceProxy))
			protected.GET("/groups/join/accept", handler.ReplyGroupJoinRequest(serviceProxy))
			protected.DELETE("/groups/:id/leave", handler.LeaveGroup(serviceProxy))

			protected.POST("/messages", handler.SendMessage(serviceProxy))
			protected.GET("/messages/offline", handler.GetOfflineMessages(serviceProxy))
			protected.PUT("/messages/:id/read", handler.MarkAsRead(serviceProxy))
			protected.PUT("/messages/read/all", handler.MarkAllAsRead(serviceProxy))
			protected.GET("/messages/unread/count", handler.GetUnreadCount(serviceProxy))
		}
	}

	// WebSocket连接
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
