package logger

import (
	"github.com/qiqiuyang/logger/model"
	"go.uber.org/zap"
	"sync"
)

var (
	once sync.Once
	log  *loggerService
)

type LoggerService interface {
	MakeLogger(model.Zap)
	GetLogger(suffix string) (*zap.Logger, bool)
	GetSugarLogger(suffix string) (*zap.SugaredLogger, bool)
	MakeDefaultLogConfig(logPath, logName, suffix string) model.Zap
}

type loggerService struct {
	loggerList     sync.Map
	logPathDefault func(logPath, logFileName string) string
}
