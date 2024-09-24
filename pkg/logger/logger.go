package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Logger struct {
	level  string
	logger *log.Logger
}

// 添加了日志等级常量，避免直接使用字符串，提高代码可读性和可维护性。
const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

// 使用log.New创建自定义的日志实例，避免使用全局的log，这样可以指定输出目的地、前缀和日志格式。
// 支持自定义输出目标io.Writer，默认为os.Stdout，方便重定向日志输出，如写入文件或其他日志收集系统。
func NewLogger(level string, out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout
	}
	return &Logger{
		level:  strings.ToLower(level),
		logger: log.New(out, "", log.LstdFlags|log.Lshortfile),
	}
}

// 使用可变参数v ...interface{}，使日志方法更加灵活，可以接受多个参数。
// 使用fmt.Sprintln格式化日志消息，保证日志输出的完整性。
func (l *Logger) logMessage(level string, v ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", strings.ToUpper(level))
	l.logger.SetPrefix(prefix)
	l.logger.Output(3, fmt.Sprintln(v...)) // 调整调用深度
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level == DEBUG {
		l.logMessage(DEBUG, v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.level == INFO || l.level == DEBUG {
		l.logMessage(INFO, v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	l.logMessage(ERROR, v...)
}

// 在Fatal方法中使用os.Exit(1)，符合log.Fatal的行为，即在输出日志后终止程序。
func (l *Logger) Fatal(v ...interface{}) {
	l.logMessage(FATAL, v...)
	os.Exit(1)
}
