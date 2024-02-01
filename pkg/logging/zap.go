package logging

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	cfg    *Config
	logger *zap.SugaredLogger
}

var once sync.Once
var zapSinLogger *zap.SugaredLogger

func newZapLogger(cfg *Config) *zapLogger {
	logger := &zapLogger{cfg: cfg}
	logger.Init()
	return logger
}

func (l *zapLogger) Init() {
	once.Do(func() {
		logger := zap.New(zapcore.NewCore(
			getEncoder(l.cfg), getWriteSyncer(l.cfg), getLogLevel(l.cfg)),
			getOptions(l.cfg)...,
		).Sugar()

		zapSinLogger = logger.With("AppName", "MyApp", "LoggerName", "Zaplog")
	})
	l.logger = zapSinLogger
}

func getEncoder(cfg *Config) zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig

	if cfg.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder

	if cfg.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	return encoder
}

func getWriteSyncer(cfg *Config) zapcore.WriteSyncer {
	return zapcore.Lock(os.Stdout)
}

var zapLogLevelMapping = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

func getLogLevel(cfg *Config) zapcore.Level {
	level, exists := zapLogLevelMapping[cfg.Level]
	if !exists {
		return zapcore.DebugLevel
	}
	return level
}

func getOptions(cfg *Config) []zap.Option {
	return []zap.Option{
		zap.AddStacktrace(zap.ErrorLevel),
		zap.AddCaller(),
	}
}

func (l *zapLogger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)

	l.logger.Debugw(msg, params...)
}

func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args)
}

func (l *zapLogger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)
	l.logger.Infow(msg, params...)
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args)
}

func (l *zapLogger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)
	l.logger.Warnw(msg, params...)
}

func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args)
}

func (l *zapLogger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)
	l.logger.Errorw(msg, params...)
}

func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args)
}

func (l *zapLogger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)
	l.logger.Fatalw(msg, params...)
}

func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args)
}

func prepareLogInfo(cat Category, sub SubCategory, extra map[ExtraKey]interface{}) []interface{} {
	if extra == nil {
		extra = make(map[ExtraKey]interface{})
	}
	extra["Category"] = cat
	extra["SubCategory"] = sub

	return logParamsToZapParams(extra)
}
