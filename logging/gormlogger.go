package logging

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

const (
	logTitle      = "[gorm] "
	messageFormat = logTitle + "%s, %s"
	traceStr      = logTitle + "[%.3fms] [rows:%v] %s"
	traceWarnStr  = logTitle + "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr   = logTitle + "%s %s\n[%.3fms] [rows:%v] %s"
)

type GormLogger struct {
	zap *ZapLogger
}

// LogMode log mode
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l *GormLogger) Info(_ context.Context, msg string, data ...interface{}) {
	l.zap.Infof(messageFormat, append([]interface{}{msg, utils.FileWithLineNum()}, data...)...)
}

// Warn print warn messages
func (l *GormLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	l.zap.Warnf(messageFormat, append([]interface{}{msg, utils.FileWithLineNum()}, data...)...)
}

// Error print error messages
func (l *GormLogger) Error(_ context.Context, msg string, data ...interface{}) {
	l.zap.Errorf(messageFormat, append([]interface{}{msg, utils.FileWithLineNum()}, data...)...)
}

// Trace print sql message
func (l *GormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	switch {
	case err != nil && (!errors.Is(err, errors.New("record not found")) || !l.zap.gormConfig.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.zap.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.zap.Errorf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.zap.gormConfig.SlowThreshold && l.zap.gormConfig.SlowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.zap.gormConfig.SlowThreshold)
		if rows == -1 {
			l.zap.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.zap.Warnf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			l.zap.Debugf(traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.zap.Debugf(traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func (l *GormLogger) ParamsFilter(_ context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.zap.gormConfig.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}
