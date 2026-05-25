// Package main 是 message 服务的入口，负责初始化 MongoDB、RabbitMQ、注册到 etcd 并启动 gRPC 服务。
package main

import (
	"IM/api/gen/message"
	"IM/pkg/config"
	"IM/pkg/discovery"
	"IM/pkg/id"
	"IM/pkg/interceptor"
	"IM/pkg/logger"
	service "IM/services/message/application"
	"IM/services/message/domain/event"
	"IM/services/message/infrastructure/mq"
	"IM/services/message/infrastructure/persistence"
	grpcserver "IM/services/message/interfaces/grpc"
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
	// 加载配置
	cfg, err := config.LoadMessageConfig(getConfigPath())
	if err != nil {
		fmt.Printf("Using default config: %v\n", err)
		cfg = config.DefaultMessageConfig()
	}
	// 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 初始化 etcd 注册中心
	registry, err := discovery.NewRegistry(cfg.Etcd.Endpoints, cfg.Etcd.DialTimeout)
	if err != nil {
		logger.Fatalw("Failed to create registry", "component", "message_cmd", "err", err)
	}
	defer registry.Close()

	// 注册服务到 etcd
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddr := cfg.Server.Host + ":" + cfg.Server.GRPCPort
	if err := registry.Register(ctx, cfg.Server.Name, serviceAddr, cfg.Etcd.TTL); err != nil {
		logger.Fatalw("Failed to register service", "component", "message_cmd", "err", err)
	}
	logger.Infow("Service registered", "component", "message_cmd", "service", cfg.Server.Name, "addr", serviceAddr)

	// 初始化 MongoDB 连接
	mongoDB, err := persistence.NewMongoDB(cfg.Database.MongoDB)
	if err != nil {
		logger.Fatalw("Failed to connect to MongoDB", "component", "message_cmd", "err", err)
	}
	defer mongoDB.Close()

	// 初始化 Redis 连接
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	// 初始化 RabbitMQ 连接
	rabbitMQ, err := mq.NewRabbitMQConnection(cfg.RabbitMQ.URL)
	if err != nil {
		logger.Fatalw("Failed to connect to RabbitMQ", "component", "message_cmd", "err", err)
	}
	defer rabbitMQ.Close()

	idGenerator := id.NewSnowflakeGenerator(3)

	// 初始化消息仓库
	messageRepo := persistence.NewMessageRepository(mongoDB)
	messageCache := persistence.NewMessageCache(redisClient)
	messageProducer := mq.NewMessageProducer(rabbitMQ, cfg.RabbitMQ.Exchange)
	// 初始化事件发布器
	eventPublisher := event.NewEventPublisher()
	messageService := service.NewMessageService(messageRepo, messageCache, messageProducer, idGenerator, eventPublisher)

	// 初始化 gRPC 服务器
	lis, err := net.Listen("tcp", ":"+cfg.Server.GRPCPort)
	if err != nil {
		logger.Fatalw("Failed to listen", "component", "message_cmd", "err", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggingUnaryInterceptor()),
	)

	messageServer := grpcserver.NewMessageServer(messageService)
	message.RegisterMessageServiceServer(grpcServer, messageServer)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(cfg.Server.Name, healthpb.HealthCheckResponse_SERVING)

	go func() {
		logger.Infow("gRPC server starting", "component", "message_cmd", "grpc_port", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalw("Failed to serve", "component", "message_cmd", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infow("Shutting down server...", "component", "message_cmd")
	grpcServer.GracefulStop()

	if err := registry.Deregister(ctx, cfg.Server.Name, serviceAddr); err != nil {
		logger.Errorw("Failed to deregister service", "component", "message_cmd", "err", err)
	}

	logger.Infow("Server stopped", "component", "message_cmd")
}

// getConfigPath 返回 message 服务使用的配置文件路径。
func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "configs/message.yaml"
}
