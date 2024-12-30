package logger

type NopLogger struct {
}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

func (zl *NopLogger) Debug(msg string, args ...Field) {
}

func (zl *NopLogger) Info(msg string, args ...Field) {
}

func (zl *NopLogger) Warn(msg string, args ...Field) {
}

func (zl *NopLogger) Error(msg string, args ...Field) {
}
