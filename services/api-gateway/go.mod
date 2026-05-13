module IM/services/api-gateway

go 1.26.1

require (
	github.com/gin-gonic/gin v1.12.0
	google.golang.org/grpc v1.81.0
	google.golang.org/grpc/credentials/insecure v1.81.0
)

replace IM => ../../..