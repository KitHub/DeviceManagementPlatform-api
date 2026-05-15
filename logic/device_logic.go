package logic

import (
	"DeviceManagementPlatform-api/dao"
	"DeviceManagementPlatform-api/wrapper"
	"context"
	"log/slog"
	"sync"
	"time"

	"xorm.io/xorm"
)

type DeviceLogic struct {
	dbEngine  *xorm.Engine
	deviceDao *dao.DeviceDAO
}

var deviceLogic *DeviceLogic
var once sync.Once

func NewDeviceLogic(ctx context.Context, dbEngine *xorm.Engine,
	deviceDao *dao.DeviceDAO) *DeviceLogic {
	once.Do(func() {
		deviceLogic = &DeviceLogic{
			dbEngine:  dbEngine,
			deviceDao: deviceDao,
		}
	})
	return deviceLogic
}

func (l *DeviceLogic) RegisterDevice(ctx context.Context,
	deviceNo string) (deviceDO *dao.DeviceDO, err error) {
	now := time.Now()
	deviceDO = &dao.DeviceDO{
		DeviceNo:     deviceNo,
		CreateTime:   now,
		UpdateTime:   now,
		RegisterTime: now,
	}
	slog.InfoContext(ctx, "register device", slog.Any("device", deviceDO))

	err = wrapper.TrsactionWrapper(ctx, l.dbEngine,
		func(session *xorm.Session) error {
			deviceDO, err = l.deviceDao.InsertDevice(ctx, session, deviceDO)
			return err
		})
	slog.InfoContext(ctx, "register device success",
		slog.Any("device", deviceDO))
	return deviceDO, nil
}

func (l *DeviceLogic) GetDeviceByNo(ctx context.Context,
	deviceNo string) (deviceDO *dao.DeviceDO, err error) {
	slog.InfoContext(ctx, "query device", slog.Any("deviceNo", deviceNo))
	deviceDO = &dao.DeviceDO{
		DeviceNo: deviceNo,
	}
	err = wrapper.TrsactionWrapper(ctx, l.dbEngine,
		func(session *xorm.Session) error {
			deviceDO, err = l.deviceDao.QueryDeviceByDeviceNo(
				ctx, session, deviceNo)
			return err
		})
	slog.InfoContext(ctx, "query device success", slog.Any("device", deviceDO))
	return deviceDO, err
}

func (l *DeviceLogic) GetDeviceById(ctx context.Context,
	deviceId int64) (deviceDO *dao.DeviceDO, err error) {
	slog.InfoContext(ctx, "query device", slog.Any("deviceId", deviceId))
	deviceDO = &dao.DeviceDO{
		Id: deviceId,
	}
	err = wrapper.TrsactionWrapper(ctx, l.dbEngine,
		func(session *xorm.Session) error {
			deviceDO, err = l.deviceDao.QueryDeviceByDeviceId(
				ctx, session, deviceId)
			return err
		})
	slog.InfoContext(ctx, "query device by id success",
		slog.Any("device", deviceDO))

	return deviceDO, err
}
