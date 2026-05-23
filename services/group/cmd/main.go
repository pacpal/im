// Package main 是 group 服务的入口，负责初始化并启动 gRPC 服务。
package main

import (
	"IM/api/gen/group"
	"IM/pkg/config"
	"IM/pkg/discovery"
	"IM/pkg/id"
	"IM/pkg/interceptor"
	"IM/pkg/logger"
	"IM/services/group/application/service"
	"IM/services/group/domain/event"
	"IM/services/group/infrastructure/persistence"
	grpcserver "IM/services/group/interfaces/grpc"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cfg, err := config.Load(getConfigPath())
	if err != nil {
		fmt.Printf("Using default config: %v\n", err)
		cfg = config.DefaultGroupConfig()
	}

	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
	}
	defer logger.Sync()

	registry, err := discovery.NewRegistry(cfg.Etcd.Endpoints, cfg.Etcd.DialTimeout)
	if err != nil {
		logger.Fatalf("Failed to create registry: %v", err)
	}
	defer registry.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddr := cfg.Server.Host + ":" + cfg.Server.GRPCPort
	if err := registry.Register(ctx, cfg.Server.Name, serviceAddr, cfg.Etcd.TTL); err != nil {
		logger.Fatalf("Failed to register service: %v", err)
	}
	logger.Infof("Service %s registered at %s", cfg.Server.Name, serviceAddr)

	db, err := persistence.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	idGenerator := id.NewSnowflakeGenerator(2)

	groupRepo := persistence.NewGroupRepository(db.GetDB())
	groupMemberRepo := persistence.NewGroupMemberRepository(db.GetDB())
	groupJoinRequestRepo := persistence.NewGroupJoinRequestRepository(db.GetDB())

	eventPublisher := event.NewEventPublisher()
	groupService := service.NewGroupService(groupRepo, groupMemberRepo, groupJoinRequestRepo, idGenerator, eventPublisher)

	lis, err := net.Listen("tcp", ":"+cfg.Server.GRPCPort)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	jwtSecret := []byte(cfg.JWT.Secret)
	skipMethods := []string{"/group.GroupService/CreateGroup", "/group.GroupService/GetGroup"}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthUnaryInterceptor(jwtSecret, skipMethods)),
	)

	groupServer := grpcserver.NewGroupServer(groupService)
	group.RegisterGroupServiceServer(grpcServer, groupServer)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(cfg.Server.Name, healthpb.HealthCheckResponse_SERVING)

	go func() {
		logger.Infof("gRPC server starting on :%s", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	grpcServer.GracefulStop()

	if err := registry.Deregister(ctx, cfg.Server.Name, serviceAddr); err != nil {
		logger.Errorf("Failed to deregister service: %v", err)
	}

	logger.Info("Server stopped")
}

// getConfigPath 返回 group 服务配置文件路径，优先使用环境变量 CONFIG_PATH。
func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "configs/group.yaml"
}
