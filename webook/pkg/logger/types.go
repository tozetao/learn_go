package logger

// LoggerV1 风格1，要求msg中包含占位符
type LoggerV1 interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Field struct {
	Key   string
	Value string
}

// LoggerV2 风格2，认为日志中的参数都是键值对。
type LoggerV2 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}
