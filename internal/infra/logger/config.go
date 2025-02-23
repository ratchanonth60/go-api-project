package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

var LoggerInstance = NewLogger()

func NewLogger() *Logger {
	stdout := zapcore.AddSync(os.Stdout)

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	productionCfg.StacktraceKey = "stack"

	jsonEncoder := zapcore.NewJSONEncoder(productionCfg)

	jsonOutCore := zapcore.NewCore(jsonEncoder, stdout, level)

	samplingCore := zapcore.NewSamplerWithOptions(
		jsonOutCore,
		time.Second, // interval
		3,           // log first 3 entries
		0,           // thereafter log zero entires within the interval
	)

	return &Logger{zap.New(samplingCore)}
}

func Info(msg string, fields ...zap.Field) {
	LoggerInstance.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	LoggerInstance.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	LoggerInstance.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	LoggerInstance.Debug(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	LoggerInstance.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	LoggerInstance.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	LoggerInstance.Fatal(msg, fields...)
}
