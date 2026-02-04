package initializers

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(serviceName string) {
	logLevel := zapcore.InfoLevel
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = zapcore.DebugLevel
	}

	writer := &lumberjack.Logger{
		Filename:   "logs/" + serviceName + ".log",
		MaxSize:    100, // MB
		MaxBackups: 7,
		MaxAge:     30, // days
		Compress:   true,
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(writer),
		logLevel,
	)

	Log = zap.New(core, zap.AddCaller())
	Log.Info("Zap 日志初始化完成！")
}
