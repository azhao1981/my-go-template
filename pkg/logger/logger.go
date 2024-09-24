package logger

import (
	"log"
	"strings"
)

type Logger struct {
	level string
}

func NewLogger(level string) *Logger {
	return &Logger{level: strings.ToLower(level)}
}

func (l *Logger) Info(msg string) {
	if l.level == "info" || l.level == "debug" {
		log.Println("[INFO]", msg)
	}
}

func (l *Logger) Debug(msg string) {
	if l.level == "debug" {
		log.Println("[DEBUG]", msg)
	}
}

func (l *Logger) Error(err error) {
	log.Println("[ERROR]", err)
}

func (l *Logger) Fatal(err error) {
	log.Fatal("[FATAL]", err)
}
