package logger

import "log"

type Logger interface {
	Log(msg string)
}

type LoggerImpl struct {
}

func NewLogger() *LoggerImpl {
	log.SetPrefix("[INFO] ")
	return &LoggerImpl{}
}

func (l *LoggerImpl) Log(msg string) {
	log.Print(msg)
}
