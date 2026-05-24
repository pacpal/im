// Package main 是 user 服务的入口，负责初始化依赖、注册到 etcd 并启动 gRPC 服务。
package main

import (
	"IM/api/gen/user"
	"IM/pkg/auth"
	"IM/pkg/config"
	"IM/pkg/discovery"
	"IM/pkg/id"
	"IM/pkg/interceptor"
	"IM/pkg/logger"
	"IM/services/user/application/service"
	"IM/services/user/domain/event"
	"IM/services/user/infrastructure/cache"
	"IM/services/user/infrastructure/persistence"
	grpcserver "IM/services/user/interfaces/grpc"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg, err := config.Load(getConfigPath())
	if err != nil {
		fmt.Printf("Using default config: %v\n", err)
		cfg = config.DefaultUserConfig()
	}

	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	registry, err := discovery.NewRegistry(cfg.Etcd.Endpoints, cfg.Etcd.DialTimeout)
	if err != nil {
		logger.Fatalw("Failed to create registry", "component", "user_cmd", "err", err)
	}
	defer registry.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddr := cfg.Server.Host + ":" + cfg.Server.GRPCPort
	if err := registry.Register(ctx, cfg.Server.Name, serviceAddr, cfg.Etcd.TTL); err != nil {
		logger.Fatalw("Failed to register service", "component", "user_cmd", "err", err)
	}
	logger.Infow("Service registered", "component", "user_cmd", "service", cfg.Server.Name, "addr", serviceAddr)

	db, err := persistence.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatalw("Failed to connect to database", "component", "user_cmd", "err", err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	idGenerator := id.NewSnowflakeGenerator(1)
	jwtUtil := auth.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.Expire)

	userRepo := persistence.NewUserRepository(db.GetDB())
	friendshipRepo := persistence.NewFriendshipRepository(db.GetDB())
	friendRequestRepo := persistence.NewFriendRequestRepository(db.GetDB())
	userCache := cache.NewUserCache(redisClient)

	eventPublisher := event.NewEventPublisher()
	userService := service.NewUserService(userRepo, friendshipRepo, friendRequestRepo, userCache, idGenerator, jwtUtil, eventPublisher)

	lis, err := net.Listen("tcp", ":"+cfg.Server.GRPCPort)
	if err != nil {
		logger.Fatalw("Failed to listen", "component", "user_cmd", "err", err)
	}

	jwtSecret := []byte(cfg.JWT.Secret)
	skipMethods := []string{"/user.UserService/Register", "/user.UserService/Login"}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthUnaryInterceptor(jwtSecret, skipMethods)),
	)

	userServer := grpcserver.NewUserServer(userService)
	user.RegisterUserServiceServer(grpcServer, userServer)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(cfg.Server.Name, healthpb.HealthCheckResponse_SERVING)

	go func() {
		logger.Infow("gRPC server starting", "component", "user_cmd", "grpc_port", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalw("Failed to serve", "component", "user_cmd", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infow("Shutting down server...", "component", "user_cmd")
	grpcServer.GracefulStop()

	if err := registry.Deregister(ctx, cfg.Server.Name, serviceAddr); err != nil {
		logger.Errorw("Failed to deregister service", "component", "user_cmd", "err", err)
	}

	logger.Infow("Server stopped", "component", "user_cmd")
}

// getConfigPath 返回 user 服务的配置文件路径，优先使用环境变量 CONFIG_PATH。
func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "configs/user.yaml"
}
