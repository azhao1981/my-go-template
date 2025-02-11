package config

import "errors"

var (
	ErrConfigFileNotFound        = errors.New("config file not found")
	ErrParseConfigFile           = errors.New("failed to parse config file")
	ErrLoadConfigFromRemote      = errors.New("failed to load config from remote")
	ErrUnmarshalConfig           = errors.New("failed to unmarshal config")
	ErrRequiredConfigItemMissing = errors.New("required config item missing")
	ErrInvalidConfigItem         = errors.New("invalid config item")
	ErrWatchConfigFailed         = errors.New("watch config failed")
)
