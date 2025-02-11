package config

import (
	"fmt"
	"os"
	"time"
)

// ExampleConfig 定义一个示例配置结构体
type ExampleConfig struct {
	Server struct { // 配置分组示例 (内嵌结构体)
		Host    string        `yaml:"host" toml:"host" json:"host" env:"SERVER_HOST" default:"localhost" comment:"服务器Host"`
		Port    int           `yaml:"port" toml:"port" json:"port" env:"SERVER_PORT" default:"8080" comment:"服务器端口" required:"true"`
		Timeout time.Duration `yaml:"timeout" toml:"timeout" json:"timeout" env:"SERVER_TIMEOUT" default:"5s" comment:"服务器超时时间"`
	} `yaml:"server" toml:"server" json:"server" comment:"服务器配置"`
	Database struct {
		Type     string `yaml:"type" toml:"type" json:"type" env:"DATABASE_TYPE" default:"mysql" comment:"数据库类型"`
		Host     string `yaml:"host" toml:"host" json:"host" env:"DATABASE_HOST" default:"localhost" comment:"数据库Host"`
		Port     int    `yaml:"port" toml:"port" json:"port" env:"DATABASE_PORT" default:"3306" comment:"数据库端口"`
		Username string `yaml:"username" toml:"username" json:"username" env:"DATABASE_USERNAME" default:"root" comment:"数据库用户名"`
		Password string `yaml:"password" toml:"password" json:"password" env:"DATABASE_PASSWORD" default:"" comment:"数据库密码"`
		DBName   string `yaml:"dbname" toml:"dbname" json:"dbname" env:"DATABASE_DBNAME" default:"mydb" comment:"数据库名"`
	} `yaml:"database" toml:"database" json:"database" comment:"数据库配置"`
	LogLevel string `yaml:"log_level" toml:"log_level" json:"log_level" env:"LOG_LEVEL" default:"info" comment:"日志级别" required:"true"`
	Metrics  struct {
		Enabled bool `yaml:"enabled" toml:"enabled" json:"enabled" env:"METRICS_ENABLED" default:"false" comment:"是否启用指标监控"`
		Port    int  `yaml:"port" toml:"port" json:"port" env:"METRICS_PORT" default:"9090" comment:"指标监控端口"`
	} `yaml:"metrics" toml:"metrics" json:"metrics" comment:"指标监控配置"`
}

// GenerateExampleConfigFiles 生成示例配置文件到当前目录
func GenerateExampleConfigFiles() error {
	configurator := NewConfigurator(&Options{}) // 使用默认 Options

	exampleConfig := &ExampleConfig{}

	yamlContent, err := configurator.GenerateExampleConfig(exampleConfig, "yaml")
	if err != nil {
		return fmt.Errorf("failed to generate yaml example config: %w", err)
	}
	if err := os.WriteFile("config.example.yaml", []byte(yamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write config.example.yaml: %w", err)
	}

	tomlContent, err := configurator.GenerateExampleConfig(exampleConfig, "toml")
	if err != nil {
		return fmt.Errorf("failed to generate toml example config: %w", err)
	}
	if err := os.WriteFile("config.example.toml", []byte(tomlContent), 0644); err != nil {
		return fmt.Errorf("failed to write config.example.toml: %w", err)
	}

	jsonContent, err := configurator.GenerateExampleConfig(exampleConfig, "json")
	if err != nil {
		return fmt.Errorf("failed to generate json example config: %w", err)
	}
	if err := os.WriteFile("config.example.json", []byte(jsonContent), 0644); err != nil {
		return fmt.Errorf("failed to write config.example.json: %w", err)
	}

	fmt.Println("Example config files (config.example.yaml, config.example.toml, config.example.json) generated successfully.")
	return nil
}
