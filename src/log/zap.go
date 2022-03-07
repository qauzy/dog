package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	ZapLog *zap.SugaredLogger // 简易版日志文件
	//logger *zap.Logger        // 这个日志强大一些, 目前还用不到

	logLevel = zap.NewAtomicLevel()
)

// InitLog 初始化日志文件
func init() {

	setLevel(zapcore.DebugLevel)

	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()) //获取编码方式

	wf := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log.txt",
		MaxSize:    128, // MB
		LocalTime:  true,
		Compress:   true,
		MaxBackups: 8, // 最多保留 n 个备份
	})

	wc := zapcore.Lock(os.Stdout) //标准输出

	file := zapcore.NewCore(encoder, wf, zapcore.DebugLevel)

	console := zapcore.NewCore(encoder, wc, zapcore.DebugLevel)

	core := zapcore.NewTee(file, console)
	ZapLog = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	return
}

func setLevel(level zapcore.Level) {
	logLevel.SetLevel(level)
}
