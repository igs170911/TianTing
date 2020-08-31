package Logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
)

var SysLog = New("Engine", zapcore.InfoLevel)

func New(service string, level zapcore.Level) *zap.SugaredLogger {
	syncWriter := getLoggerSync(service)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, syncWriter, level)
	logger := zap.New(core)
	sugaredLogger := logger.Sugar()
	defer sugaredLogger.Sync()
	return sugaredLogger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)

	// 如果要給 Json 的話
	//encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//return zapcore.NewJSONEncoder(encoderConfig)
}

func getLoggerSync(service string) zapcore.WriteSyncer {
	var ws []zapcore.WriteSyncer
	var logRoot string
	if runtime.GOOS == "linux" {
		logRoot = "/var/log"
	} else {
		logRoot = "./log"
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", logRoot, service),
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	ws = append(ws, zapcore.AddSync(lumberJackLogger))
	ws = append(ws, zapcore.AddSync(os.Stdout))
	return zapcore.NewMultiWriteSyncer(ws...)
}
