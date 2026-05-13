package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig 从 YAML 文件加载配置到 out 结构体
func LoadConfig(path string, out interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}
