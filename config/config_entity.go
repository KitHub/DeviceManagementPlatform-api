package config

type LogConfigEntity struct {
	LogLevel   string `yaml:"log_level"`   // 日志级别（如：debug、info、warn、error）
	Filename   string `yaml:"filename"`    // 日志文件的位置
	MaxSize    int    `yaml:"max_size"`    // 文件最大尺寸（以MB为单位）
	MaxBackups int    `yaml:"max_backups"` // 保留的最大旧文件数量
	MaxAge     int    `yaml:"max_age"`     // 保留旧文件的最大天数
	Compress   bool   `yaml:"compress"`    // 是否压缩/归档旧文件
	LocalTime  bool   `yaml:"local_time"`  // 使用本地时间创建时间戳
}

type ConfigEntity struct {
	LogConfig *LogConfigEntity `yaml:"log"`
}
