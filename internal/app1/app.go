package app1

import (
	"myGoTemplate/internal/app1/handler"
	"myGoTemplate/pkg/logger"
	"net/http"
	"strconv"
)

type Config struct {
	Port   int
	Logger *logger.Logger
}

func Run(cfg *Config, logger *logger.Logger) error {
	// 设置路由和处理器
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handler.HelloHandler)

	// 启动 HTTP 服务器
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: mux,
	}

	logger.Info("Starting app1 on port " + strconv.Itoa(cfg.Port))

	return server.ListenAndServe()
}
