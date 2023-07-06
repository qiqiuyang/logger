package model

type Zap struct {
	Level         string `json:"level"`          // 级别
	Format        string `json:"format"`         // 输出
	Suffix        string `json:"suffix"`         // 日志后缀
	ShowLine      bool   `json:"show-line"`      // 显示行
	EncodeLevel   string ` json:"encode-level"`  // 编码级
	StacktraceKey string `json:"stacktrace-key"` // 栈名
	LogInConsole  bool   `json:"log-in-console"` // 输出控制台
	FilePath      string `json:"file-path"`      // 日志输出路径
	MaxSize       int    // 在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups    int    // 保留旧文件的最大个数
	MaxAge        int    // 保留旧文件的最大天数
	Compress      bool   // 是否压缩/归档旧文件
	CallerKey     string // 设置后才可以打印文件名和行号

}
