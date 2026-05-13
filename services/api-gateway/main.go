package main

import (
	"IM/services/api-gateway/config"
	"IM/services/api-gateway/handler"
	"IM/services/api-gateway/proxy"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetDefaultConfig()

	serviceProxy, err := proxy.NewServiceProxy(cfg)
	if err != nil {
		log.Fatalf("Failed to create service proxy: %v", err)
	}

	gatewayHandler := handler.NewGatewayHandler(serviceProxy, cfg)

	router := gin.Default()
	gatewayHandler.RegisterRoutes(router)

	log.Printf("API Gateway starting on :%s", cfg.Gateway.Port)
	if err := router.Run(":" + cfg.Gateway.Port); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}