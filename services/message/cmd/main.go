package main

import (
	"IM/api/gen/message"
	"IM/pkg/config"
	"IM/pkg/discovery"
	"IM/pkg/id"
	"IM/pkg/interceptor"
	"IM/services/message/application/service"
	"IM/services/message/domain/event"
	"IM/services/message/infrastructure/mq"
	"IM/services/message/infrastructure/persistence"
	grpcserver "IM/services/message/interfaces/grpc"
	"context"
	"log"
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
	cfg, err := config.LoadMessageConfig(getConfigPath())
	if err != nil {
		log.Printf("Using default config: %v", err)
		cfg = config.DefaultMessageConfig()
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

	mongoDB, err := persistence.NewMongoDB(cfg.Database.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	rabbitMQ, err := mq.NewRabbitMQConnection(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	idGenerator := id.NewSnowflakeGenerator(3)

	messageRepo := persistence.NewMessageRepository(mongoDB)
	messageCache := persistence.NewMessageCache(redisClient)
	messageProducer := mq.NewMessageProducer(rabbitMQ, cfg.RabbitMQ.Exchange)

	eventPublisher := event.NewEventPublisher()
	messageService := service.NewMessageService(messageRepo, messageCache, messageProducer, idGenerator, eventPublisher)

	lis, err := net.Listen("tcp", ":"+cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	jwtSecret := []byte(cfg.JWT.Secret)
	skipMethods := []string{}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthUnaryInterceptor(jwtSecret, skipMethods)),
	)

	messageServer := grpcserver.NewMessageServer(messageService)
	message.RegisterMessageServiceServer(grpcServer, messageServer)

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
	return "configs/message.yaml"
}
