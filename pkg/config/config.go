// Package config 提供应用配置结构、默认配置以及从 YAML 文件加载配置的辅助函数。
package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 是应用的主配置结构，包含服务、Etcd、数据库、Redis、JWT 和日志配置。
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Etcd     EtcdConfig     `yaml:"etcd"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
}

// ServerConfig 表示单个服务的网络配置（名称、HTTP 端口、gRPC 端口和主机）。
type ServerConfig struct {
	Name     string `yaml:"name"`
	HTTPPort string `yaml:"http_port"`
	GRPCPort string `yaml:"grpc_port"`
	Host     string `yaml:"host"`
}

// EtcdConfig 表示 etcd 的连接选项和租约 TTL 等配置。
type EtcdConfig struct {
	Endpoints   []string      `yaml:"endpoints"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	TTL         int64         `yaml:"ttl"`
}

// DatabaseConfig 表示关系型数据库连接相关配置。
type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	MaxOpen  int    `yaml:"max_open"`
	MaxIdle  int    `yaml:"max_idle"`
}

// RedisConfig 表示 Redis 连接配置。
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// JWTConfig 表示 JWT 的密钥和过期时长配置。
type JWTConfig struct {
	Secret string        `yaml:"secret"`
	Expire time.Duration `yaml:"expire"`
}

// LogConfig 表示日志级别与格式配置。
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// GatewayConfig 是 Gateway 服务使用的配置结构，包含服务端、etcd、各下游服务和其他资源配置。
type GatewayConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Etcd     EtcdConfig     `yaml:"etcd"`
	Services ServicesConfig `yaml:"services"`
	Redis    RedisConfig    `yaml:"redis"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
}

// ServicesConfig 包含各后端微服务的端点配置。
type ServicesConfig struct {
	User    ServiceEndpointConfig `yaml:"user"`
	Group   ServiceEndpointConfig `yaml:"group"`
	Message ServiceEndpointConfig `yaml:"message"`
}

// ServiceEndpointConfig 表示单个服务的名称与 gRPC 端口信息。
type ServiceEndpointConfig struct {
	Name     string `yaml:"name"`
	GRPCPort string `yaml:"grpc_port"`
}

// MessageConfig 是消息服务的配置结构，包含数据库、RabbitMQ、WebSocket 等配置。
type MessageConfig struct {
	Server    ServerConfig    `yaml:"server"`
	Etcd      EtcdConfig      `yaml:"etcd"`
	Database  MessageDBConfig `yaml:"database"`
	Redis     RedisConfig     `yaml:"redis"`
	RabbitMQ  RabbitMQConfig  `yaml:"rabbitmq"`
	WebSocket WebSocketConfig `yaml:"websocket"`
	JWT       JWTConfig       `yaml:"jwt"`
	Log       LogConfig       `yaml:"log"`
}

// MessageDBConfig 包含消息服务使用的 Postgres 与 MongoDB 的配置。
type MessageDBConfig struct {
	Postgres DatabaseConfig `yaml:"postgres"`
	MongoDB  MongoDBConfig  `yaml:"mongodb"`
}

// MongoDBConfig 表示 MongoDB 的连接和集合配置。
type MongoDBConfig struct {
	URI        string `yaml:"uri"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

// RabbitMQConfig 表示 RabbitMQ 连接与交换机相关配置。
type RabbitMQConfig struct {
	URL         string `yaml:"url"`
	Exchange    string `yaml:"exchange"`
	QueuePrefix string `yaml:"queue_prefix"`
}

// WebSocketConfig 定义 WebSocket 的缓冲、心跳与消息大小配置。
type WebSocketConfig struct {
	ReadBufferSize  int           `yaml:"read_buffer_size"`
	WriteBufferSize int           `yaml:"write_buffer_size"`
	PingPeriod      time.Duration `yaml:"ping_period"`
	PongWait        time.Duration `yaml:"pong_wait"`
	MaxMessageSize  int           `yaml:"max_message_size"`
}

// Load 从指定路径读取 YAML 文件并解析为 Config 结构。
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadGatewayConfig 从指定路径读取 Gateway 专用的 YAML 配置并解析为 GatewayConfig。
func LoadGatewayConfig(path string) (*GatewayConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg GatewayConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadMessageConfig 从指定路径读取 Message 服务的 YAML 配置并解析为 MessageConfig。
func LoadMessageConfig(path string) (*MessageConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg MessageConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// DefaultUserConfig 返回一个用于本地开发的 user 服务默认配置。
func DefaultUserConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Name:     "user-service",
			HTTPPort: "8081",
			GRPCPort: "50051",
			Host:     "0.0.0.0",
		},
		Etcd: EtcdConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			TTL:         30,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "im_db",
			SSLMode:  "disable",
			MaxOpen:  100,
			MaxIdle:  10,
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: "6379",
		},
		RabbitMQ: RabbitMQConfig{
			URL:         "amqp://guest:guest@localhost:5672/",
			Exchange:    "im_events",
			QueuePrefix: "im_user_",
		},
		JWT: JWTConfig{
			Secret: "your-secret-key",
			Expire: 24 * time.Hour,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}
}

// DefaultGroupConfig 返回 group 服务的默认开发配置。
func DefaultGroupConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Name:     "group-service",
			HTTPPort: "8082",
			GRPCPort: "50052",
			Host:     "0.0.0.0",
		},
		Etcd: EtcdConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			TTL:         30,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "im_db",
			SSLMode:  "disable",
			MaxOpen:  100,
			MaxIdle:  10,
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: "6379",
		},
		RabbitMQ: RabbitMQConfig{
			URL:         "amqp://guest:guest@localhost:5672/",
			Exchange:    "im_exchange",
			QueuePrefix: "im_",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}
}

// DefaultMessageConfig 返回 message 服务的默认开发配置。
func DefaultMessageConfig() *MessageConfig {
	return &MessageConfig{
		Server: ServerConfig{
			Name:     "message-service",
			HTTPPort: "8083",
			GRPCPort: "50053",
			Host:     "0.0.0.0",
		},
		Etcd: EtcdConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
			TTL:         30,
		},
		Database: MessageDBConfig{
			Postgres: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "postgres",
				DBName:   "im_db",
				SSLMode:  "disable",
				MaxOpen:  100,
				MaxIdle:  10,
			},
			MongoDB: MongoDBConfig{
				URI:        "mongodb://localhost:27017",
				Database:   "im_messages",
				Collection: "messages",
			},
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: "6379",
		},
		RabbitMQ: RabbitMQConfig{
			URL:         "amqp://guest:guest@localhost:5672/",
			Exchange:    "im_exchange",
			QueuePrefix: "im_",
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			PingPeriod:      30 * time.Second,
			PongWait:        60 * time.Second,
			MaxMessageSize:  65536,
		},
		JWT: JWTConfig{
			Secret: "your-secret-key",
			Expire: 24 * time.Hour,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}
}

// DefaultGatewayConfig 返回 gateway 的默认配置（用于本地/开发环境）。
func DefaultGatewayConfig() *GatewayConfig {
	return &GatewayConfig{
		Server: ServerConfig{
			Name:     "gateway",
			HTTPPort: "8080",
			Host:     "0.0.0.0",
		},
		Etcd: EtcdConfig{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		},
		Services: ServicesConfig{
			User: ServiceEndpointConfig{
				Name:     "user-service",
				GRPCPort: "50051",
			},
			Group: ServiceEndpointConfig{
				Name:     "group-service",
				GRPCPort: "50052",
			},
			Message: ServiceEndpointConfig{
				Name:     "message-service",
				GRPCPort: "50053",
			},
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: "6379",
		},
		RabbitMQ: RabbitMQConfig{
			URL:         "amqp://guest:guest@localhost:5672/",
			Exchange:    "im_exchange",
			QueuePrefix: "im_",
		},
		JWT: JWTConfig{
			Secret: "your-secret-key",
			Expire: 24 * time.Hour,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}
}
