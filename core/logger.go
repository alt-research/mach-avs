package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
)

type ZapLogger struct {
	logger *zap.Logger
}

var _ sdklogging.Logger = (*ZapLogger)(nil)

func NewZapLogger(env sdklogging.LogLevel) (sdklogging.Logger, error) {
	config := zap.NewProductionConfig()
	if env == sdklogging.Development {
		config = zap.NewDevelopmentConfig()
	}

	config.DisableStacktrace = true
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	return &ZapLogger{
		logger: logger,
	}, nil
}

func (z *ZapLogger) Debug(msg string, tags ...any) {
	z.logger.Sugar().Debugw(msg, tags...)
}

func (z *ZapLogger) Info(msg string, tags ...any) {
	z.logger.Sugar().Infow(msg, tags...)
}

func (z *ZapLogger) Warn(msg string, tags ...any) {
	z.logger.Sugar().Warnw(msg, tags...)
}

func (z *ZapLogger) Error(msg string, tags ...any) {
	z.logger.Sugar().Errorw(msg, tags...)
}

func (z *ZapLogger) Fatal(msg string, tags ...any) {
	z.logger.Sugar().Fatalw(msg, tags...)
}

func (z *ZapLogger) Debugf(template string, args ...interface{}) {
	z.logger.Sugar().Debugf(template, args...)
}

func (z *ZapLogger) Infof(template string, args ...interface{}) {
	z.logger.Sugar().Infof(template, args...)
}

func (z *ZapLogger) Warnf(template string, args ...interface{}) {
	z.logger.Sugar().Warnf(template, args...)
}

func (z *ZapLogger) Errorf(template string, args ...interface{}) {
	z.logger.Sugar().Errorf(template, args...)
}

func (z *ZapLogger) Fatalf(template string, args ...interface{}) {
	z.logger.Sugar().Fatalf(template, args...)
}

func (z *ZapLogger) With(tags ...any) sdklogging.Logger {
	return &ZapLogger{
		logger: z.logger.Sugar().With(tags...).Desugar(),
	}
}
