package logging

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	gormLogger "gorm.io/gorm/logger"
	"os"
	"time"
	"vet-clinic/config"
)

type ZapLogger struct {
	zap        *zap.SugaredLogger
	gormConfig *gormLogger.Config
}

func NewLogger(zap *zap.SugaredLogger, gorm *gormLogger.Config) *ZapLogger {
	return &ZapLogger{zap: zap, gormConfig: gorm}
}

func Init(cfg *config.Config) (l *ZapLogger) {
	logger, err := build(cfg)
	if err != nil {
		fmt.Printf("Failed to compose zap logger : %s", err)
		os.Exit(config.ErrExitStatus)
	}
	sugar := logger.Sugar()
	l = &ZapLogger{
		zap:        sugar,
		gormConfig: buildGormConfig(cfg),
	}
	l.Infof("Success to read zap logger configuration")
	_ = logger.Sync()
	return
}

// Zap returns zap.SugaredLogger
func (l *ZapLogger) Zap() *zap.SugaredLogger {
	return l.zap
}

// Gorm wraps the Logger to provide an implementation of the GORM logger interface
func (l *ZapLogger) Gorm() gormLogger.Interface {
	return &GormLogger{zap: l}
}
func (l *ZapLogger) Debugf(msg string, args ...interface{}) {
	l.zap.Debugf(msg, args...)
}
func (l *ZapLogger) Infof(msg string, args ...interface{}) {
	l.zap.Infof(msg, args...)
}
func (l *ZapLogger) Warnf(msg string, args ...interface{}) {
	l.zap.Warnf(msg, args...)
}
func (l *ZapLogger) Errorf(msg string, args ...interface{}) {
	l.zap.Errorf(msg, args...)
}
func (l *ZapLogger) Panicf(msg string, args ...interface{}) {
	l.zap.Panicf(msg, args...)
}
func (l *ZapLogger) Fatalf(msg string, args ...interface{}) {
	l.zap.Fatalf(msg, args...)
}

func build(cfg *config.Config) (*zap.Logger, error) {
	var zapCfg = cfg.Logger.ZapConfig
	enc, _ := newEncoder(zapCfg)
	writer, errWriter := openWriters(cfg)

	if zapCfg.Level == (zap.AtomicLevel{}) {
		return nil, errors.New("missing Level")
	}

	log := zap.New(zapcore.NewCore(enc, writer, zapCfg.Level), buildOptions(zapCfg, errWriter)...)
	return log, nil
}

func newEncoder(cfg zap.Config) (zapcore.Encoder, error) {
	switch cfg.Encoding {
	case "console":
		return zapcore.NewConsoleEncoder(cfg.EncoderConfig), nil
	case "json":
		return zapcore.NewJSONEncoder(cfg.EncoderConfig), nil
	}
	return nil, errors.New("failed to set encoder")
}

func openWriters(cfg *config.Config) (zapcore.WriteSyncer, zapcore.WriteSyncer) {
	writer := open(cfg.Logger.ZapConfig.OutputPaths, &cfg.Logger.LogRotate)
	errWriter := open(cfg.Logger.ZapConfig.ErrorOutputPaths, &cfg.Logger.LogRotate)
	return writer, errWriter
}

func open(paths []string, rotateCfg *lumberjack.Logger) zapcore.WriteSyncer {
	writers := make([]zapcore.WriteSyncer, 0, len(paths))
	for _, path := range paths {
		writer := newWriter(path, rotateCfg)
		writers = append(writers, writer)
	}
	writer := zap.CombineWriteSyncers(writers...)
	return writer
}

func newWriter(path string, rotateCfg *lumberjack.Logger) zapcore.WriteSyncer {
	switch path {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	}
	sink := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   path,
			MaxSize:    rotateCfg.MaxSize,
			MaxBackups: rotateCfg.MaxBackups,
			MaxAge:     rotateCfg.MaxAge,
			Compress:   rotateCfg.Compress,
		},
	)
	return sink
}

func buildOptions(cfg zap.Config, errWriter zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errWriter)}
	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if cfg.Development {
		stackLevel = zap.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}
	return opts
}

func buildGormConfig(cfg *config.Config) *gormLogger.Config {
	slowThreshold := 200 * time.Millisecond
	if cfg.Logger.GormConfig.SlowThreshold > 0 {
		slowThreshold = cfg.Logger.GormConfig.SlowThreshold
	}
	return &gormLogger.Config{
		SlowThreshold:             slowThreshold,
		IgnoreRecordNotFoundError: cfg.Logger.GormConfig.IgnoreRecordNotFoundError,
		ParameterizedQueries:      cfg.Logger.GormConfig.ParameterizedQueries,
	}
}
