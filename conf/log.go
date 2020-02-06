package conf

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var SugarLogger *zap.SugaredLogger
// 初始化log
func InitLog() {
	level := GetLogLevel()
	writeSyncer, _ := os.Create(Config.Log.Path)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core)
	SugarLogger = logger.Sugar()
	defer SugarLogger.Sync()
}

func GetLogLevel() zapcore.Level {
	level := Config.Log.Level
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	}
	return zapcore.DebugLevel
}