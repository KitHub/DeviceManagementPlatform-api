package servicecontext

import (
	"DeviceManagementPlatform-api/config"
	"DeviceManagementPlatform-api/logic"
	"context"
	"log/slog"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
	"xorm.io/xorm"
)

var serviceLogger *slog.Logger

type ServiceContext struct {
	Logger      *slog.Logger
	DB          *xorm.Engine
	DeviceLogic *logic.DeviceLogic
}

var serviceCtx *ServiceContext

func InitServiceContext(ctx context.Context, configEntity config.ConfigEntity) (
	serviceCtx *ServiceContext, err error) {
	logger, err := initLog(ctx, configEntity.LogConfig)
	if err != nil {
		slog.ErrorContext(ctx, "init log failed", slog.Any("error", err))
		return nil, err
	}
	db, err := initDB(ctx, configEntity.DBConfig, logger)
	if err != nil {
		slog.ErrorContext(ctx, "init db failed", slog.Any("error", err))
		return nil, err
	}

	deviceLogic := logic.NewDeviceLogic(ctx, db)

	serviceCtx = &ServiceContext{
		Logger:      logger,
		DB:          db,
		DeviceLogic: deviceLogic,
	}
	return serviceCtx, nil
}

func initLog(ctx context.Context, logConfig *config.LogConfigEntity) (
	*slog.Logger, error) {
	log := &lumberjack.Logger{
		Filename:   logConfig.Filename,   // 日志文件路径
		MaxSize:    logConfig.MaxSize,    // 每个日志文件的最大大小（以MB为单位）
		MaxBackups: logConfig.MaxBackups, // 保留旧文件的最大数量
		MaxAge:     logConfig.MaxAge,     // 保留旧文件的最大天数
		Compress:   logConfig.Compress,   // 是否压缩旧文件
		LocalTime:  logConfig.LocalTime,  // 是否使用本地时间戳
	}
	serviceLogger := slog.New(slog.NewTextHandler(log, nil))
	slog.SetDefault(serviceLogger)
	return serviceLogger, nil
}

func initDB(ctx context.Context, dbConfig *config.DBConfigEntity,
	logger *slog.Logger) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(dbConfig.DriverName, dbConfig.DataSourceName)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to initialize database connection",
			slog.Any("error", err))
		return nil, err
	}

	if logger != nil {
		engine.SetLogger(logger)
	}
	engine.ShowSQL(dbConfig.ShowSQL)
	engine.SetMaxIdleConns(dbConfig.MaxIdleConns)
	engine.SetMaxOpenConns(dbConfig.MaxOpenConns)
	engine.SetConnMaxLifetime(
		time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

	err = engine.PingContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to ping database", slog.Any("error", err))
		return nil, err
	}

	return engine, nil
}
