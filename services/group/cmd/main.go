package main

import (
	"IM/api/gen/group"
	"IM/pkg/config"
	"IM/pkg/discovery"
	"IM/pkg/id"
	"IM/pkg/interceptor"
	"IM/services/group/application/service"
	"IM/services/group/domain/event"
	"IM/services/group/infrastructure/persistence"
	grpcserver "IM/services/group/interfaces/grpc"
	"context"
	"log"
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
		log.Printf("Using default config: %v", err)
		cfg = config.DefaultGroupConfig()
	}

	registry, err := discovery.NewRegistry(cfg.Etcd.Endpoints, cfg.Etcd.DialTimeout)
	if err != nil {
		log.Fatalf("Failed to create registry: %v", err)
	}
	defer registry.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceAddr := cfg.Server.Host + ":" + cfg.Server.GRPCPort
	if err := registry.Register(ctx, cfg.Server.Name, serviceAddr, cfg.Etcd.TTL); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	log.Printf("Service %s registered at %s", cfg.Server.Name, serviceAddr)

	db, err := persistence.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
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
		log.Fatalf("Failed to listen: %v", err)
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
		log.Printf("gRPC server starting on :%s", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()

	if err := registry.Deregister(ctx, cfg.Server.Name, serviceAddr); err != nil {
		log.Printf("Failed to deregister service: %v", err)
	}

	log.Println("Server stopped")
}

func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "configs/group.yaml"
}
