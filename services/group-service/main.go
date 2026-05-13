package main

import (
	"IM/services/group-service/application/service"
	"IM/services/group-service/config"
	grpcserver "IM/services/group-service/interface/grpc"
	"IM/services/group-service/interface/handler"
	"IM/server/repository/postgres"
	"context"
	"log"
	"net"

	"IM/api/gen/group"
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

	groupRepo := postgres.NewGroupRepo(dbHandler.GetDB())
	groupMemberRepo := postgres.NewGroupMemberRepo(dbHandler.GetDB())
	groupJoinRequestRepo := postgres.NewGroupJoinRequestRepo(dbHandler.GetDB())
	userRepo := postgres.NewUserRepo(dbHandler.GetDB())

	groupSvc := service.NewGroupApplicationService(groupRepo, groupMemberRepo, groupJoinRequestRepo, userRepo)

	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		group.RegisterGroupServiceServer(grpcServer, grpcserver.NewGroupGrpcServer(groupSvc))
		log.Printf("gRPC server starting on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	router := gin.Default()
	groupHandler := handler.NewGroupHttpHandler(groupSvc)

	apiGroup := router.Group("/api/v1")
	groupHandler.RegisterRoutes(apiGroup, groupSvc)

	log.Printf("HTTP server starting on :8082")
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