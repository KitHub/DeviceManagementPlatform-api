package logic

import (
	"DeviceManagementPlatform-api/db"
	"context"
	"sync"
	"time"

	"xorm.io/xorm"
)

type DeviceLogic struct {
	dbEngine  *xorm.Engine
	deviceDao *db.DeviceDAO
}

var deviceLogic *DeviceLogic
var once sync.Once

func NewDeviceLogic(ctx context.Context, dbEngine *xorm.Engine) *DeviceLogic {
	once.Do(func() {
		deviceLogic = &DeviceLogic{
			dbEngine:  dbEngine,
			deviceDao: db.NewDeviceDAO(dbEngine),
		}
	})
	return deviceLogic
}

func (l *DeviceLogic) RegisterDevice(ctx context.Context,
	deviceNo string) (*db.DeviceDO, error) {
	now := time.Now()
	deviceDO := &db.DeviceDO{
		DeviceNo:     deviceNo,
		CreateTime:   now,
		UpdateTime:   now,
		RegisterTime: now,
	}

	session := l.dbEngine.NewSession()
	_, err := l.deviceDao.InsertDevice(ctx, session, deviceDO)
	if err != nil {
		return nil, err
	}
	return deviceDO, nil
}

func (l *DeviceLogic) GetDeviceByNo(ctx context.Context,
	deviceNo string) (*db.DeviceDO, error) {
	deviceDO := &db.DeviceDO{
		DeviceNo: deviceNo,
	}
	has, err := l.dbEngine.Get(deviceDO)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return deviceDO, nil
}
