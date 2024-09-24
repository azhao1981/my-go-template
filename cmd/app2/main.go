package main

import (
	"myGoTemplate/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	err := logger.Init()
	if err != nil {
		panic(err)
	}
	defer logger.ZipLogger.Sync()

	logger.ZipLogger.Info("应用程序已启动")

	// 示例日志
	logger.ZipLogger.Debug("这是一个调试日志")
	logger.ZipLogger.Warn("这是一个警告日志")
	logger.ZipLogger.Error("这是一个错误日志")

	// 使用带字段的日志
	logger.ZipLogger.Info("用户登录", zap.String("username", "admin"), zap.Int("userID", 123))
}
