package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"testing"
)

func TestSlice(t *testing.T) {
	str := "1234567890"
	fmt.Printf("%v\n", str[0:2])
	fmt.Printf("%v\n", str[0:0])
}

func TestViper(t *testing.T) {
	InitConfig()

	fmt.Println(viper.AllSettings())
}

func TestLog(t *testing.T) {
	logger, err := zap.NewProduction()

	if err != nil {
		panic(err)
	}

	logger.Debug("This is a debug.")
	logger.Info("This is a info.")
	logger.Warn("This is a warn.")

	//encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//
	//fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	//level := zapcore.DebugLevel
	//
	//logFile, err := os.OpenFile("./log-test-zap.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 06666)
	//if err != nil {
	//	panic(err)
	//}
	//writer := zapcore.AddSync(logFile)
	//
	//logger := zap.New(zapcore.NewCore(fileEncoder, writer, level),
	//	zap.AddCaller(),
	//	zap.AddStacktrace(zapcore.ErrorLevel))
	//defer logger.Sync()
	//
	//url := "http://www.test.com"
	//logger.Info("write log to file",
	//	zap.String("url", url),
	//	zap.Int("attempt", 3))

}
