package servicecontext

import (
	"DeviceManagementPlatform-api/config"
	"DeviceManagementPlatform-api/dao"
	"DeviceManagementPlatform-api/logic"
	"DeviceManagementPlatform-api/service"
	"context"
	"log/slog"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
	"xorm.io/xorm"
)

var serviceLogger *slog.Logger

type ShutdownCallback func(ctx context.Context) error

var shutdownCallbacks []ShutdownCallback

type ServiceContext struct {
	Logger      *slog.Logger
	DB          *xorm.Engine
	DeviceDao   *dao.DeviceDAO
	DeviceLogic *logic.DeviceLogic
	ApiService  *service.ApiService
}

var gServiceCtx *ServiceContext
var once sync.Once

func RegisterShutdownCallback(callback ShutdownCallback) {
	shutdownCallbacks = append(shutdownCallbacks, callback)
}

func GetShutdownCallbacks() []ShutdownCallback {
	return shutdownCallbacks
}

func InitServiceContext(ctx context.Context, configEntity *config.ConfigEntity) (
	serviceCtx *ServiceContext, err error) {
	slog.InfoContext(ctx, "init service context")
	once.Do(func() {
		logger, innerErr := initLog(ctx, configEntity.LogConfig)
		if innerErr != nil {
			slog.ErrorContext(ctx, "init log failed", slog.Any("error", innerErr))
			err = innerErr
			return
		}
		db, innerErr := initDB(ctx, configEntity.DBConfig, logger)
		if innerErr != nil {
			slog.ErrorContext(ctx, "init db failed", slog.Any("error", innerErr))
			err = innerErr
			return
		}

		deviceDao := dao.NewDeviceDAO(ctx)
		deviceLogic := logic.NewDeviceLogic(ctx, db, deviceDao)
		apiService := service.NewApiService(ctx, deviceLogic)

		gServiceCtx = &ServiceContext{
			Logger:      logger,
			DB:          db,
			DeviceDao:   deviceDao,
			DeviceLogic: deviceLogic,
			ApiService:  apiService,
		}
	})
	if err != nil {
		slog.ErrorContext(ctx, "init service context failed", slog.Any("error", err))
		return nil, err

	}
	slog.InfoContext(ctx, "init service context done")
	return gServiceCtx, err

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

	engine.SetMaxIdleConns(dbConfig.MaxIdleConns)
	engine.SetMaxOpenConns(dbConfig.MaxOpenConns)
	engine.SetConnMaxLifetime(
		time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

	err = engine.PingContext(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to ping database", slog.Any("error", err))
		return nil, err
	}

	RegisterShutdownCallback(func(ctx context.Context) error {
		return engine.Close()
	})

	return engine, nil
}
