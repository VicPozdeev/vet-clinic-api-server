package logging

import (
	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
)

type Logger interface {
	Zap() *zap.SugaredLogger
	Gorm() gormLogger.Interface
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Panicf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}
