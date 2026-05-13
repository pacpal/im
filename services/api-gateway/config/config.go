package config

type Config struct {
	Gateway  GatewayConfig
	Services ServicesConfig
}

type GatewayConfig struct {
	Port string
	Host string
}

type ServicesConfig struct {
	UserService    ServiceConfig
	GroupService   ServiceConfig
	MessageService ServiceConfig
}

type ServiceConfig struct {
	Host string
	Port string
	GRPC GRPCConfig
}

type GRPCConfig struct {
	Port string
}

func GetDefaultConfig() *Config {
	return &Config{
		Gateway: GatewayConfig{
			Port: "8080",
			Host: "0.0.0.0",
		},
		Services: ServicesConfig{
			UserService: ServiceConfig{
				Host: "localhost",
				Port: "8081",
				GRPC: GRPCConfig{
					Port: "50051",
				},
			},
			GroupService: ServiceConfig{
				Host: "localhost",
				Port: "8082",
				GRPC: GRPCConfig{
					Port: "50052",
				},
			},
			MessageService: ServiceConfig{
				Host: "localhost",
				Port: "8083",
				GRPC: GRPCConfig{
					Port: "50053",
				},
			},
		},
	}
}