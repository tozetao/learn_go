package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.Logger
}

func NewLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: logger}
}

func (zl *ZapLogger) Debug(msg string, args ...Field) {
	zl.logger.Debug(msg, zl.toArgs(args)...)
}

func (zl *ZapLogger) Info(msg string, args ...Field) {
	zl.logger.Info(msg, zl.toArgs(args)...)
}

func (zl *ZapLogger) Warn(msg string, args ...Field) {
	zl.logger.Warn(msg, zl.toArgs(args)...)
}

func (zl *ZapLogger) Error(msg string, args ...Field) {
	zl.logger.Error(msg, zl.toArgs(args)...)
}

func (zl *ZapLogger) toArgs(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}
	return res
}
