package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Etcd     EtcdConfig     `yaml:"etcd"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Name     string `yaml:"name"`
	HTTPPort string `yaml:"http_port"`
	GRPCPort string `yaml:"grpc_port"`
	Host     string `yaml:"host"`
}

type EtcdConfig struct {
	Endpoints    []string      `yaml:"endpoints"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	TTL          int64         `yaml:"ttl"`
}

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

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	Secret string        `yaml:"secret"`
	Expire time.Duration `yaml:"expire"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type GatewayConfig struct {
	Server   ServerConfig        `yaml:"server"`
	Etcd     EtcdConfig          `yaml:"etcd"`
	Services ServicesConfig      `yaml:"services"`
	Redis    RedisConfig         `yaml:"redis"`
	JWT      JWTConfig           `yaml:"jwt"`
	Log      LogConfig           `yaml:"log"`
}

type ServicesConfig struct {
	User    ServiceEndpointConfig `yaml:"user"`
	Group   ServiceEndpointConfig `yaml:"group"`
	Message ServiceEndpointConfig `yaml:"message"`
}

type ServiceEndpointConfig struct {
	Name     string `yaml:"name"`
	GRPCPort string `yaml:"grpc_port"`
}

type MessageConfig struct {
	Server    ServerConfig     `yaml:"server"`
	Etcd      EtcdConfig       `yaml:"etcd"`
	Database  MessageDBConfig  `yaml:"database"`
	Redis     RedisConfig      `yaml:"redis"`
	RabbitMQ  RabbitMQConfig   `yaml:"rabbitmq"`
	WebSocket WebSocketConfig  `yaml:"websocket"`
	Log       LogConfig        `yaml:"log"`
}

type MessageDBConfig struct {
	Postgres DatabaseConfig `yaml:"postgres"`
	MongoDB  MongoDBConfig  `yaml:"mongodb"`
}

type MongoDBConfig struct {
	URI        string `yaml:"uri"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

type RabbitMQConfig struct {
	URL          string `yaml:"url"`
	Exchange     string `yaml:"exchange"`
	QueuePrefix  string `yaml:"queue_prefix"`
}

type WebSocketConfig struct {
	ReadBufferSize  int           `yaml:"read_buffer_size"`
	WriteBufferSize int           `yaml:"write_buffer_size"`
	PingPeriod      time.Duration `yaml:"ping_period"`
	PongWait        time.Duration `yaml:"pong_wait"`
	MaxMessageSize  int           `yaml:"max_message_size"`
}

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
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}
}

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
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
	}
}

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
