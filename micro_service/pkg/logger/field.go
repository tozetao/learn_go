package logger

type Field struct {
	Key   string
	Value any
}

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}

func String(key string, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Int(key string, val int) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}
func Int32(key string, val int32) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}
func Int64(key string, val int64) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}

func Bool(key string, val bool) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}
