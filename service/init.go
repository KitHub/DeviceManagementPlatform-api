package service

import (
	"DeviceManagementPlatform-api/config"
	"context"
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitService(ctx context.Context, configEntity config.ConfigEntity) error {
	InitLog(ctx, configEntity.LogConfig)
	return nil
}

func InitLog(ctx context.Context, logConfig *config.LogConfigEntity) {
	log := &lumberjack.Logger{
		Filename:   logConfig.Filename,   // 日志文件路径
		MaxSize:    logConfig.MaxSize,    // 每个日志文件的最大大小（以MB为单位）
		MaxBackups: logConfig.MaxBackups, // 保留旧文件的最大数量
		MaxAge:     logConfig.MaxAge,     // 保留旧文件的最大天数
		Compress:   logConfig.Compress,   // 是否压缩旧文件
		LocalTime:  logConfig.LocalTime,  // 是否使用本地时间戳
	}
	textLogger := slog.New(slog.NewTextHandler(log, nil))
	slog.SetDefault(textLogger)
}
