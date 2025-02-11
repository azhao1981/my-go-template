package config

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// Loader 接口定义了配置加载器的行为
type Loader interface {
	Load(ctx context.Context, config interface{}) error
}

// FileLoader 从文件加载配置
type FileLoader struct {
	paths     []string
	filenames []string
}

func NewFileLoader(paths []string, filenames []string) *FileLoader {
	return &FileLoader{paths: paths, filenames: filenames}
}

func (l *FileLoader) Load(ctx context.Context, config interface{}) error {
	v := viper.New()
	v.SetConfigType("yaml") // 默认 YAML, 自动识别 yml, yaml

	for _, filename := range l.filenames {
		for _, path := range l.paths {
			filepath := filepath.Join(path, filename)
			if _, err := os.Stat(filepath); err == nil {
				v.SetConfigFile(filepath)
				if err := v.ReadInConfig(); err != nil {
					if _, ok := err.(viper.ConfigFileNotFoundError); ok {
						continue // 文件不存在，尝试下一个路径或文件名
					} else {
						return fmt.Errorf("failed to read config file: %w", err) // 其他读取错误
					}
				}
				if err := v.Unmarshal(config); err != nil {
					return fmt.Errorf("failed to unmarshal config from file: %w", err)
				}
				return nil // 成功加载并Unmarshal后直接返回
			}
		}
	}

	return ErrConfigFileNotFound // 所有文件路径都找不到文件
}

// EnvLoader 从环境变量加载配置
type EnvLoader struct {
	prefix string
}

func NewEnvLoader(prefix string) *EnvLoader {
	return &EnvLoader{prefix: prefix}
}

func (l *EnvLoader) Load(ctx context.Context, config interface{}) error {
	v := viper.New()
	v.SetEnvPrefix(l.prefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 替换 . 为 _，例如 db.host 变为 DB_HOST

	// 使用 gotenv 加载 .env 文件 (可选)
	gotenv.Load()

	if err := v.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config from env: %w", err)
	}
	return nil
}

// RemoteLoader 从远程 HTTP JSON 加载配置
type RemoteLoader struct {
	endpoint string
	interval time.Duration
	client   *http.Client
}

func NewRemoteLoader(endpoint string, interval time.Duration) *RemoteLoader {
	client := &http.Client{Timeout: 10 * time.Second} // 设置超时
	return &RemoteLoader{endpoint: endpoint, interval: interval, client: client}
}

func (l *RemoteLoader) Load(ctx context.Context, config interface{}) error {
	v := viper.New()
	resp, err := l.client.Get(l.endpoint)
	if err != nil {
		return fmt.Errorf("failed to get remote config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get remote config, status code: %d", resp.StatusCode)
	}

	v.SetConfigType("json") // 远程配置通常是 JSON
	if err := v.ReadConfig(resp.Body); err != nil {
		return fmt.Errorf("failed to read remote config body: %w", err)
	}

	if err := v.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config from remote: %w", err)
	}
	return nil
}
