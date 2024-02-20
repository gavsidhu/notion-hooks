package logging

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.Formatter = &logrus.JSONFormatter{}

	Logger.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/notion-hooks.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})
}
