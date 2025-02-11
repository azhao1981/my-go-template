# Golang Configuration Management Module

This module provides a robust and flexible configuration management solution for Golang applications, based on best practices like Domain-Driven Design (DDD), SOLID principles, and KISS principle. It leverages popular libraries like `viper`, `cobra`, `gotenv`, and `go-defaults`.

## Features

*   **Struct-Based Configuration Definition:** Define your application configuration using Go structs with tags for metadata (type, default value, environment variable binding, comments, required).
*   **Configuration Grouping:** Organize configurations using nested structs.
*   **Multiple Configuration Sources:**
    *   **Environment Variables:** Load configuration from environment variables with prefix support.
    *   **YAML/TOML Files:** Load configuration from local YAML or TOML files.
    *   **Remote HTTP JSON:** Load configuration from remote HTTP JSON endpoints with polling.
*   **Viper Integration:**  Underlyingly powered by `viper`, providing access to all of Viper's features.
*   **Configuration Merging:**  Merges configurations from different sources with environment variables having the highest priority.
*   **Automatic Default Values:**  Populate default values from struct tags using `go-defaults`.
*   **Required Field Validation:**  Enforce required configuration fields using struct tags, preventing application startup with missing critical configurations.
*   **Configuration Watching and Hot Reloading:** Automatically reload configuration on a time interval and trigger `OnChange` callbacks for specific configuration items.
*   **Configuration Export:** Export the current configuration to YAML, TOML, or JSON formats.
*   **Example Configuration Generation:** Generate example configuration files (YAML, TOML, JSON) with default values to help users get started.

## Getting Started

### 1. Define your configuration struct

Create a Go struct in your application to define your configuration. Use struct tags to specify metadata:

```go
package main

import "time"

type Config struct {
	Server struct {
		Host    string        `yaml:"host" env:"SERVER_HOST" default:"localhost" comment:"Server host"`
		Port    int           `yaml:"port" env:"SERVER_PORT" default:"8080" required:"true" comment:"Server port"`
		Timeout time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT" default:"5s" comment:"Server timeout"`
	} `yaml:"server" comment:"Server configuration"`
	Database struct {
		Host     string `yaml:"host" env:"DATABASE_HOST" default:"localhost" comment:"Database host"`
		Port     int    `yaml:"port" env:"DATABASE_PORT" default:"3306" comment:"Database port"`
		Username string `yaml:"username" env:"DATABASE_USERNAME" default:"user" comment:"Database username"`
		Password string `yaml:"password" env:"DATABASE_PASSWORD" default:"password" comment:"Database password"`
		DBName   string `yaml:"dbname" env:"DATABASE_DBNAME" default:"mydb" comment:"Database name"`
	} `yaml:"database" comment:"Database configuration"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL" default:"info" comment:"Log level"`
}
```

### 2. Load configuration in your `main.go`

```go
package main

import (
	"context"
	"log"
	"time"

	"your-module/pkg/config" // 替换为你的模块路径
)

func main() {
	opts := &config.Options{
		ConfigFiles: []string{"config.yaml", "app-config.toml"}, // 配置文件列表
		FilePaths:   []string{"./config", "."},                 // 配置文件搜索路径
		EnvPrefix:   "APP",                                     // 环境变量前缀
		RemoteConfig: config.RemoteOptions{                      // 远程配置
			Provider: "http",
			Endpoint: "http://localhost:8081/config.json",
			Interval: 10 * time.Minute,
		},
		WatchInterval: 5 * time.Minute, // 配置监控间隔
	}

	cfg := &Config{}
	configurator := config.NewConfigurator(opts)

	if err := configurator.LoadConfig(context.Background(), cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Start config watching (optional)
	configurator.StartWatch(context.Background(), cfg)

	// Example: Register OnChange callback
	configurator.OnChange("log_level", func() {
		log.Println("Log level changed to:", configurator.GetViper().GetString("log_level"))
		// Perform actions on config change, e.g., reconfigure logging
	})


	// Access configurations using cfg struct or viper
	log.Println("Server Host:", cfg.Server.Host)
	log.Println("Database Host (via viper):", configurator.GetViper().GetString("database.host"))


	// ... your application logic ...


	// Generate example config files (for development or documentation)
	if err := config.GenerateExampleConfigFiles(); err != nil {
		log.Println("Failed to generate example config files:", err)
	}

	// Export current config (for debugging or auditing)
	exportedConfig, err := configurator.ExportConfig("yaml")
	if err != nil {
		log.Println("Failed to export config:", err)
	} else {
		log.Println("Current Config (YAML):\n", exportedConfig)
	}


	// Keep the application running (for watching example)
	select {}
}
```

### 3. Create `config.yaml` (or other config files)

Place your `config.yaml` (or other configured files) in the specified paths (e.g., `./config/config.yaml` or `./config.yaml`).

```yaml
# config.yaml
server:
  port: 8080 # Override default port
database:
  host: config-db-host # Example config file value
log_level: warn # Example config file value
metrics:
  enabled: true
```

### 4. Run your application

```bash
go run main.go
```

You can also set environment variables (e.g., `APP_SERVER_HOST=env-host APP_LOG_LEVEL=debug`) to override configurations from files and defaults.

## Configuration Options (`config.yml`)

The configuration module itself can be configured using a `config.yml` file (default name).  See `pkg/config/options.go` for available options.  Example `config.yml`:

```yaml
config_files:
  - app-config.yaml
  - custom-config.toml
file_paths:
  - /opt/app/config
  - ./
remote_config:
  provider: http
  endpoint: http://config-server/config.json
  interval: 60s
env_prefix: MY_APP
watch_interval: 30s
```

## Contributing

... (Contribution guidelines) ...

## License

... (License information) ... 