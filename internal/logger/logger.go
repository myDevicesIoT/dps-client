package logger

import (
	log "github.com/sirupsen/logrus"
)

type LevelLogger struct {
	level log.Level
}

func New(level log.Level) *LevelLogger {
	return &LevelLogger{
		level: level,
	}
}

func (logger *LevelLogger) Println(args ...interface{}) {
	log.StandardLogger().Logln(logger.level, args...)
}

func (logger *LevelLogger) Printf(format string, args ...interface{}) {
	log.StandardLogger().Logf(logger.level, format, args...)
}
