package logger

import (
	"fmt"
	"github.com/qiqiuyang/logger/model"
	"go.uber.org/zap"
	"os"
	"os/user"
	"path"
	"runtime"
	"sync"
)

const defaultLogFileName = "default.log"

func NewLoggerService(logPathFunc func(logPath, logFileName string) string) LoggerService {
	once.Do(func() {
		log = &loggerService{
			loggerList:     sync.Map{},
			logPathDefault: getLogPath,
		}
		if logPathFunc != nil {
			log.logPathDefault = logPathFunc
		}
	})

	return log
}

// MakeLogger 根据配置生成并存储logger
func (l *loggerService) MakeLogger(config model.Zap) {
	l.loggerList.Store(config.Suffix, Zap(config))
}

// GetLogger 普通打印，不原生支持sprintf类格式化语句
func (l *loggerService) GetLogger(suffix string) (*zap.Logger, bool) {
	if value, ok := l.loggerList.Load(suffix); ok {
		return value.(*zap.Logger), ok
	}
	return nil, false
}

// GetSugarLogger 高级打印，支持sprintf类格式化语句
func (l *loggerService) GetSugarLogger(suffix string) (*zap.SugaredLogger, bool) {
	if value, ok := l.GetLogger(suffix); ok {
		return value.Sugar(), ok
	}
	return nil, false
}

// MakeDefaultLogConfig 生成默认配置
func (l *loggerService) MakeDefaultLogConfig(logPath, logName, suffix string) model.Zap {
	return model.Zap{
		Level:         "info",                             // 级别
		Format:        "console",                          // 输出
		Suffix:        fmt.Sprintf("[%s]", suffix),        // 日志后缀
		ShowLine:      false,                              // 显示行
		EncodeLevel:   "CapitalLevelEncoder",              // 编码级，默认使用 LowercaseLevelEncoder，大写不带颜色(在文件里带颜色无法显示)
		StacktraceKey: "stacktrace",                       // 栈名
		LogInConsole:  true,                               // 输出控制台
		FilePath:      l.logPathDefault(logPath, logName), // 日志文件的位置
		MaxSize:       10,                                 // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups:    200,                                // 保留旧文件的最大个数
		MaxAge:        30,                                 // 保留旧文件的最大天数
		Compress:      false,                              // 是否压缩/归档旧文件
		CallerKey:     "",
	}
}

// 获取文件存储路径, 默认mosn日志存放路径+输入日志文件名
func getLogPath(logPath, logFileName string) string {
	var logFolder string
	if logFolder != "" {
		logFolder = logPath
	} else {
		if u, err := user.Current(); err != nil {
			logFolder = "/home/admin/logs/envExporter"
		} else if runtime.GOOS == "darwin" {
			logFolder = path.Join(u.HomeDir, "logs/envExporter")
		} else if runtime.GOOS == "windows" {
			logFolder = path.Join(u.HomeDir, "logs/envExporter")
		} else {
			logFolder = "/home/admin/logs/envExporter"
		}
	}
	_, err := os.Stat(logFolder)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(logFolder, 0755)
	}

	_ = os.Chmod(logFolder, 0755)

	return logFolder + "/" + logFileName
}
