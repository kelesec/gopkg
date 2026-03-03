package logger

type Option func(*Config)

// WithLevel 设置日志级别
func WithLevel(level Level) Option {
	return func(c *Config) {
		c.level = level
	}
}

// WithLogDir 设置日志目录
func WithLogDir(dir string) Option {
	return func(c *Config) {
		c.logDir = dir
	}
}

// WithLogFile 设置日志文件名
func WithLogFile(filename string) Option {
	return func(c *Config) {
		c.logFile = filename
	}
}

// WithConsoleLog 启用/禁用控制台输出
func WithConsoleLog(enable bool) Option {
	return func(c *Config) {
		c.allowConsoleLog = enable
	}
}

// WithFileLog 启用/禁用文件输出
func WithFileLog(enable bool) Option {
	return func(c *Config) {
		c.allowFileLog = enable
	}
}

// WithMaxSize 设置单个日志文件最大尺寸（MB）
func WithMaxSize(size int) Option {
	return func(c *Config) {
		c.maxSize = size
	}
}

// WithMaxBackups 设置最大备份数
func WithMaxBackups(backups int) Option {
	return func(c *Config) {
		c.maxBackups = backups
	}
}

// WithMaxAge 设置日志保留天数
func WithMaxAge(age int) Option {
	return func(c *Config) {
		c.maxAge = age
	}
}

// WithCompress 设置是否压缩
func WithCompress(compress bool) Option {
	return func(c *Config) {
		c.compress = compress
	}
}

// WithJSONFormat 设置是否使用JSON格式（文件日志）
func WithJSONFormat(enable bool) Option {
	return func(c *Config) {
		c.jsonFormat = enable
	}
}

// WithCaller 设置是否显示调用者信息
func WithCaller(enable bool) Option {
	return func(c *Config) {
		c.caller = enable
	}
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		level:           InfoLevel,
		logDir:          "",
		logFile:         "",
		allowConsoleLog: true,
		allowFileLog:    false,
		maxSize:         100,
		maxBackups:      30,
		maxAge:          30,
		compress:        true,
		jsonFormat:      false,
		caller:          true,
	}
}
