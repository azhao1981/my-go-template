package config

import "time"

// Options 定义了配置模块的选项，从 config.yml 加载
type Options struct {
	ConfigFiles   []string      `yaml:"config_files" toml:"config_files" json:"config_files" default:"config.yml" comment:"配置文件列表，支持yml, toml, json"`
	FilePaths     []string      `yaml:"file_paths" toml:"file_paths" json:"file_paths" default:"." comment:"配置文件搜索路径"`
	RemoteConfig  RemoteOptions `yaml:"remote_config" toml:"remote_config" json:"remote_config" comment:"远程配置"`
	EnvPrefix     string        `yaml:"env_prefix" toml:"env_prefix" json:"env_prefix" default:"" comment:"环境变量前缀"`
	WatchInterval time.Duration `yaml:"watch_interval" toml:"watch_interval" json:"watch_interval" default:"1m" comment:"配置监控间隔"`
}

// RemoteOptions 定义远程配置选项
type RemoteOptions struct {
	Provider  string        `yaml:"provider" toml:"provider" json:"provider" comment:"远程配置提供者，例如 http"`
	Endpoint  string        `yaml:"endpoint" toml:"endpoint" json:"endpoint" comment:"远程配置端点"`
	SecretKey string        `yaml:"secret_key" toml:"secret_key" json:"secret_key" comment:"远程配置密钥"`
	Interval  time.Duration `yaml:"interval" toml:"interval" json:"interval" default:"5m" comment:"远程配置更新间隔"`
}
