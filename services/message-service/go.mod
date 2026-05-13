module IM/services/message-service

go 1.26.1

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/gorilla/websocket v1.5.3
	google.golang.org/grpc v1.81.0
)

replace IM => ../../..