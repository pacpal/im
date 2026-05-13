package config

import "os"

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host string
	Port string
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := parseYAML(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseYAML(data []byte, out interface{}) error {
	return nil
}

func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8083",
			Host: "0.0.0.0",
		},
		DB: DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "im_db",
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: "6379",
		},
	}
}