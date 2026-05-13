package main

import (
	"IM/services/user-service/application/service"
	"IM/services/user-service/config"
	"IM/services/user-service/interface/handler"
	grpcserver "IM/services/user-service/interface/grpc"
	"IM/server/repository/postgres"
	"context"
	"log"
	"net"
	"time"

	"IM/api/gen/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.GetDefaultConfig()

	dbHandler, err := postgres.NewDBHandler(postgres.DBConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbHandler.Close()

	userRepo := postgres.NewUserRepo(dbHandler.GetDB())
	friendshipRepo := postgres.NewFriendshipRepo(dbHandler.GetDB())
	friendRequestRepo := postgres.NewFriendRequestRepo(dbHandler.GetDB())

	userSvc := service.NewUserApplicationService(userRepo, friendshipRepo, friendRequestRepo)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		user.RegisterUserServiceServer(grpcServer, grpcserver.NewUserGrpcServer(userSvc))
		log.Printf("gRPC server starting on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	router := gin.Default()
	userHandler := handler.NewUserHttpHandler(userSvc)

	apiGroup := router.Group("/api/v1")
	userHandler.RegisterRoutes(apiGroup, userSvc)

	log.Printf("HTTP server starting on :8081")
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}

type DBHandler struct {
	db interface{}
}

func NewDBHandler(cfg postgres.DBConfig) (*DBHandler, error) {
	return &DBHandler{}, nil
}

func (h *DBHandler) Close() error {
	return nil
}
