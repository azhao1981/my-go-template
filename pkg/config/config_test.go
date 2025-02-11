package config

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 定义一个测试用的配置结构体
type TestConfig struct {
	Server struct {
		Host string `yaml:"host" toml:"host" json:"host" env:"TEST_SERVER_HOST" default:"test-host"`
		Port int    `yaml:"port" toml:"port" json:"port" env:"TEST_SERVER_PORT" default:"9999" required:"true"`
	} `yaml:"server" toml:"server" json:"server"`
	LogLevel string `yaml:"log_level" toml:"log_level" json:"log_level" env:"TEST_LOG_LEVEL" default:"debug"`
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	configurator := NewConfigurator(&Options{})
	testConfig := &TestConfig{}

	err := configurator.LoadConfig(context.Background(), testConfig)
	assert.NoError(t, err)
	assert.Equal(t, "test-host", testConfig.Server.Host)
	assert.Equal(t, 9999, testConfig.Server.Port)
	assert.Equal(t, "debug", testConfig.LogLevel)
}

func TestLoadConfig_EnvValues(t *testing.T) {
	os.Setenv("TEST_SERVER_HOST", "env-host")
	os.Setenv("TEST_SERVER_PORT", "8080") // 环境变量是字符串
	os.Setenv("TEST_LOG_LEVEL", "info")
	defer os.Clearenv() // 清理环境变量

	configurator := NewConfigurator(&Options{EnvPrefix: "TEST"}) // 设置环境变量前缀
	testConfig := &TestConfig{}

	err := configurator.LoadConfig(context.Background(), testConfig)
	assert.NoError(t, err)
	assert.Equal(t, "env-host", testConfig.Server.Host)
	assert.Equal(t, 8080, testConfig.Server.Port) // 自动转换为 int
	assert.Equal(t, "info", testConfig.LogLevel)
}

func TestLoadConfig_RequiredFieldMissing(t *testing.T) {
	os.Clearenv() // 确保没有设置 TEST_SERVER_PORT 环境变量
	configurator := NewConfigurator(&Options{EnvPrefix: "TEST"})
	testConfig := &TestConfig{}

	err := configurator.LoadConfig(context.Background(), testConfig)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrRequiredConfigItemMissing)
	assert.Contains(t, err.Error(), "Server.Port") // 错误信息包含字段名
}

func TestConfigurator_OnChange(t *testing.T) {
	configurator := NewConfigurator(&Options{WatchInterval: 1 * time.Second}) // 启用监控
	testConfig := &TestConfig{}
	configurator.LoadConfig(context.Background(), testConfig) // 首次加载

	changeCount := 0
	configurator.OnChange("log_level", func() {
		changeCount++
		newLogLevel := configurator.GetViper().GetString("log_level")
		assert.Equal(t, "warn", newLogLevel, "OnChange callback should reflect new log level")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 监控3秒
	defer cancel()
	configurator.StartWatch(ctx, testConfig)

	// 模拟配置变更 (实际场景可能是文件或远程配置更新)
	v := configurator.GetViper()
	v.Set("log_level", "warn")  // 直接修改 Viper 中的值来模拟配置变更 (测试场景)
	time.Sleep(2 * time.Second) // 等待监控间隔触发

	assert.GreaterOrEqual(t, changeCount, 1, "OnChange callback should be triggered")
}

func TestConfigurator_ExportConfig(t *testing.T) {
	configurator := NewConfigurator(&Options{})
	testConfig := &TestConfig{}
	configurator.LoadConfig(context.Background(), testConfig)

	yamlContent, err := configurator.ExportConfig("yaml")
	assert.NoError(t, err)
	assert.NotEmpty(t, yamlContent)
	assert.Contains(t, yamlContent, "host: test-host") // 包含配置项和值

	jsonContent, err := configurator.ExportConfig("json")
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonContent)
	assert.Contains(t, jsonContent, "\"host\":\"test-host\"")

	tomlContent, err := configurator.ExportConfig("toml")
	assert.NoError(t, err)
	assert.NotEmpty(t, tomlContent)
	assert.Contains(t, tomlContent, "host = \"test-host\"")

	_, err = configurator.ExportConfig("invalid")
	assert.Error(t, err, "should return error for invalid format")
}

func TestConfigurator_GenerateExampleConfig(t *testing.T) {
	configurator := NewConfigurator(&Options{})
	exampleConfig := &ExampleConfig{}

	yamlContent, err := configurator.GenerateExampleConfig(exampleConfig, "yaml")
	assert.NoError(t, err)
	assert.NotEmpty(t, yamlContent)
	assert.Contains(t, yamlContent, "host: localhost") // 包含默认值和注释

	jsonContent, err := configurator.GenerateExampleConfig(exampleConfig, "json")
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonContent)
	assert.Contains(t, jsonContent, "\"host\": \"localhost\"")

	tomlContent, err := configurator.GenerateExampleConfig(exampleConfig, "toml")
	assert.NoError(t, err)
	assert.NotEmpty(t, tomlContent)
	assert.Contains(t, tomlContent, "host = \"localhost\"")

	_, err = configurator.GenerateExampleConfig(exampleConfig, "invalid")
	assert.Error(t, err, "should return error for invalid format")
}
