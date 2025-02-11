package config

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/mcuadros/go-defaults"
	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Configurator 是配置管理器的核心结构体
type Configurator struct {
	options       *Options
	viper         *viper.Viper
	onChangeFuncs map[string][]func() // 配置项变更回调函数
	mu            sync.RWMutex
}

// NewConfigurator 创建一个新的 Configurator 实例
func NewConfigurator(opts *Options) *Configurator {
	return &Configurator{
		options:       opts,
		viper:         viper.New(),
		onChangeFuncs: make(map[string][]func()),
	}
}

// LoadConfig 加载配置
func (c *Configurator) LoadConfig(ctx context.Context, config interface{}) error {
	options := c.options

	// 1. 加载默认值 (通过 struct tag)
	defaults.SetDefaults(config)

	// 2. 加载配置文件
	fileLoader := NewFileLoader(options.FilePaths, options.ConfigFiles)
	if err := fileLoader.Load(ctx, config); err != nil && err != ErrConfigFileNotFound {
		log.Printf("failed to load config from file: %v", err) // 文件加载失败，但不是文件不存在，记录日志
	}

	// 3. 加载远程配置
	if options.RemoteConfig.Endpoint != "" {
		remoteLoader := NewRemoteLoader(options.RemoteConfig.Endpoint, options.RemoteConfig.Interval)
		if err := remoteLoader.Load(ctx, config); err != nil {
			log.Printf("failed to load config from remote: %v", err)
		}
	}

	// 4. 加载环境变量 (优先级最高)
	envLoader := NewEnvLoader(options.EnvPrefix)
	if err := envLoader.Load(ctx, config); err != nil {
		log.Printf("failed to load config from env: %v", err)
	}

	// 5. Unmarshal to viper and validate required fields
	if err := c.unmarshalAndValidate(config); err != nil {
		return err
	}

	log.Println("config loaded successfully")
	return nil
}

func (c *Configurator) unmarshalAndValidate(config interface{}) error {
	v := c.viper

	// 将配置项绑定到 Viper
	val := reflect.ValueOf(config).Elem() // 获取结构体的值
	typ := val.Type()                     // 获取结构体的类型

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		tag := field.Tag

		// 处理内嵌结构体
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			embeddedVal := val.Field(i)
			embeddedType := embeddedVal.Type()
			for j := 0; j < embeddedType.NumField(); j++ {
				embeddedField := embeddedType.Field(j)
				embeddedFieldName := embeddedField.Name
				embeddedTag := embeddedField.Tag
				envName := embeddedTag.Get("env")
				if envName == "" {
					envName = strings.ToUpper(c.options.EnvPrefix + "_" + fieldName + "_" + embeddedFieldName) // 默认环境变量名
				}
				viperFieldName := strings.ToLower(fieldName + "." + embeddedFieldName) // viper 中的配置名
				v.BindEnv(viperFieldName, envName)
				v.SetDefault(viperFieldName, embeddedTag.Get("default"))
			}
			continue // 跳过内嵌结构体自身的处理
		}

		envName := tag.Get("env")
		if envName == "" {
			envName = strings.ToUpper(c.options.EnvPrefix + "_" + fieldName) // 默认环境变量名
		}
		viperFieldName := strings.ToLower(fieldName) // viper 中的配置名

		v.BindEnv(viperFieldName, envName)
		v.SetDefault(viperFieldName, tag.Get("default"))

		required := tag.Get("required")
		if required == "true" {
			if !v.IsSet(viperFieldName) || v.Get(viperFieldName) == "" {
				return fmt.Errorf("%w: %s", ErrRequiredConfigItemMissing, fieldName)
			}
		}
	}

	if err := v.Unmarshal(config); err != nil { // 再次 Unmarshal，确保 Viper 的值被应用到 config 结构体
		return fmt.Errorf("failed to unmarshal config to struct after viper setup: %w", err)
	}

	c.viper = v // 更新 Configurator 中的 viper 实例
	return nil
}

// GetViper 返回底层的 viper 实例，方便用户直接使用 viper 的功能
func (c *Configurator) GetViper() *viper.Viper {
	return c.viper
}

// OnChange 注册配置项变更回调函数
func (c *Configurator) OnChange(item string, fn func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onChangeFuncs[item] = append(c.onChangeFuncs[item], fn)
}

// StartWatch 启动配置监控
func (c *Configurator) StartWatch(ctx context.Context, config interface{}) error {
	interval := c.options.WatchInterval
	if interval <= 0 {
		return nil // 不监控
	}

	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := c.reloadConfig(ctx, config); err != nil {
					log.Printf("failed to reload config: %v", err)
				}
			}
		}
	}()
	return nil
}

// reloadConfig 重新加载配置并触发 OnChange 回调
func (c *Configurator) reloadConfig(ctx context.Context, config interface{}) error {
	oldValues := c.getCurrentConfigValues(config) // 获取旧的配置值

	if err := c.LoadConfig(ctx, config); err != nil { // 重新加载配置
		return fmt.Errorf("failed to reload config: %w", err)
	}

	newValues := c.getCurrentConfigValues(config) // 获取新的配置值
	c.triggerOnChangeEvents(oldValues, newValues) // 触发 OnChange 事件
	return nil
}

// getCurrentConfigValues 获取当前配置项的值 (用于比较新旧值)
func (c *Configurator) getCurrentConfigValues(config interface{}) map[string]interface{} {
	values := make(map[string]interface{})
	v := c.viper

	val := reflect.ValueOf(config).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		if field.Type.Kind() == reflect.Struct && field.Anonymous { // 处理内嵌结构体
			embeddedVal := val.Field(i)
			embeddedType := embeddedVal.Type()
			for j := 0; j < embeddedType.NumField(); j++ {
				embeddedField := embeddedType.Field(j)
				embeddedFieldName := embeddedField.Name
				viperFieldName := strings.ToLower(fieldName + "." + embeddedFieldName)
				values[viperFieldName] = v.Get(viperFieldName)
			}
			continue
		}

		viperFieldName := strings.ToLower(fieldName)
		values[viperFieldName] = v.Get(viperFieldName)
	}
	return values
}

// triggerOnChangeEvents 触发配置项变更事件
func (c *Configurator) triggerOnChangeEvents(oldValues, newValues map[string]interface{}) {
	c.mu.RLock() // 读锁
	defer c.mu.RUnlock()

	for item, newV := range newValues {
		oldV, ok := oldValues[item]
		if ok && !reflect.DeepEqual(oldV, newV) { // 配置项值发生变化
			if fns, exist := c.onChangeFuncs[item]; exist {
				for _, fn := range fns {
					go fn() // 异步执行回调函数
				}
			}
		}
	}
}

// ExportConfig 导出当前配置到指定格式 (yaml, toml, json)
func (c *Configurator) ExportConfig(format string) (string, error) {
	v := c.viper
	settings := v.AllSettings() // 获取所有配置

	switch strings.ToLower(format) {
	case "yaml":
		out, err := yaml.Marshal(settings)
		if err != nil {
			return "", fmt.Errorf("failed to marshal config to yaml: %w", err)
		}
		return string(out), nil
	case "toml":
		out, err := toml.Marshal(settings) // 使用 toml.Marshal
		if err != nil {
			return "", fmt.Errorf("failed to marshal config to toml: %w", err)
		}
		return string(out), nil
	case "json":
		out, err := json.Marshal(settings) // 使用 json.Marshal
		if err != nil {
			return "", fmt.Errorf("failed to marshal config to json: %w", err)
		}
		return string(out), nil
	default:
		return "", fmt.Errorf("unsupported config format: %s, supported formats are yaml, toml, json", format)
	}
}

// GenerateExampleConfig 生成示例配置文件内容 (yaml, toml, json)
func (c *Configurator) GenerateExampleConfig(config interface{}, format string) (string, error) {
	v := viper.New()
	defaults.SetDefaults(config) // 设置默认值

	// 将默认值写入 viper
	val := reflect.ValueOf(config).Elem()
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		tag := field.Tag
		defaultValue := tag.Get("default")
		viperFieldName := strings.ToLower(fieldName)
		v.SetDefault(viperFieldName, defaultValue)
	}

	settings := v.AllSettings()

	switch strings.ToLower(format) {
	case "yaml":
		v.SetConfigType("yaml")
		out, err := yaml.Marshal(settings)
		if err != nil {
			return "", fmt.Errorf("failed to marshal example config to yaml: %w", err)
		}
		return "# Example Config (YAML)\n# Please modify the values as needed\n\n" + string(out), nil
	case "toml":
		out, err := toml.Marshal(settings) // 使用 toml.Marshal
		if err != nil {
			return "", fmt.Errorf("failed to marshal config to toml: %w", err)
		}
		return "# Example Config (TOML)\n# Please modify the values as needed\n\n" + string(out), nil
	case "json":
		out, err := json.Marshal(settings) // 使用 json.Marshal
		if err != nil {
			return "", fmt.Errorf("failed to marshal config to json: %w", err)
		}
		return "// Example Config (JSON)\n// Please modify the values as needed\n\n" + string(out), nil
	default:
		return "", fmt.Errorf("unsupported config format: %s, supported formats are yaml, toml, json", format)
	}
}
