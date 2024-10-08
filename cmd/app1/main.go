package main

import (
	"myGoTemplate/internal/app1"
	"myGoTemplate/pkg/config"
	"myGoTemplate/pkg/logger"
	"os"
)

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// 初始化日志
	logger := logger.NewLogger("info", os.Stdout)

	// 启动应用程序
	if err := app1.Run(cfg, logger); err != nil {
		logger.Fatal(err)
	}
}
