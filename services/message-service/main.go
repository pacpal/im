package main

import (
	"IM/services/message-service/application/service"
	"IM/services/message-service/config"
	grpcserver "IM/services/message-service/interface/grpc"
	"IM/services/message-service/interface/handler"
	"IM/server/repository/postgres"
	"log"
	"net"

	"IM/api/gen/message"
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

	messageRepo := postgres.NewMessageRepo(dbHandler.GetDB())
	friendshipRepo := postgres.NewFriendshipRepo(dbHandler.GetDB())
	groupMemberRepo := postgres.NewGroupMemberRepo(dbHandler.GetDB())

	msgSvc := service.NewMessageApplicationService(messageRepo, friendshipRepo, groupMemberRepo)
	msgSvc.StartHub()

	go func() {
		lis, err := net.Listen("tcp", ":50053")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		message.RegisterMessageServiceServer(grpcServer, grpcserver.NewMessageGrpcServer(msgSvc))
		log.Printf("gRPC server starting on :50053")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	router := gin.Default()
	msgHandler := handler.NewMessageHttpHandler(msgSvc)

	apiGroup := router.Group("/api/v1")
	msgHandler.RegisterRoutes(apiGroup, msgSvc)

	log.Printf("HTTP server starting on :8083")
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