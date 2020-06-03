package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap *zap.Logger

type level struct {
	// 最低写入等级
	lowestLevel zapcore.Level
	// 最高写入等级
	highestLevel zapcore.Level
}

func (l level) Enabled(lv zapcore.Level) bool {
	return l.lowestLevel <= lv && lv <= l.highestLevel
}

// 切割，分类
func init() {
	var coreTree zapcore.Core
	infoCfg := DefaultWriteConfig()
	infoCfg.Filename = "./log/info_logger.log"
	infoEnc := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	infoLevel := level{
		lowestLevel:  zapcore.DebugLevel,
		highestLevel: zapcore.InfoLevel,
	}
	infoCore := zapcore.NewCore(infoEnc, NewWrite(infoCfg), infoLevel)

	errCfg := DefaultWriteConfig()
	errCfg.Filename = "./log/err_logger.log"
	errEnc := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	errLevel := level{
		lowestLevel:  zapcore.WarnLevel,
		highestLevel: zapcore.FatalLevel,
	}
	errCore := zapcore.NewCore(errEnc, NewWrite(errCfg), errLevel)

	coreTree = zapcore.NewTee(infoCore, errCore)
	Zap = zap.New(coreTree, zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel))
}
