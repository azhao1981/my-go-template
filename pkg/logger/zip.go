package logger

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ZipLogger *zap.Logger

// Init 初始化日志器
func Init() error {
	// 设置默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.output", "stdout")

	// 读取配置文件
	viper.SetConfigName("config/app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // 可以根据需要修改配置文件路径

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 获取日志配置
	levelStr := viper.GetString("log.level")
	output := viper.GetString("log.output")

	// 解析日志级别
	var logLevel zapcore.Level
	if err := logLevel.UnmarshalText([]byte(levelStr)); err != nil {
		logLevel = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// 配置输出
	var writeSyncer zapcore.WriteSyncer
	switch output {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	default:
		file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writeSyncer = zapcore.AddSync(file)
	}

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		writeSyncer,
		logLevel,
	)

	// 构建日志器
	ZipLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}
