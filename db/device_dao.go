package db

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type DeviceDAO struct {
}

var deviceDAO *DeviceDAO
var once sync.Once

func NewDeviceDAO(dbEngine *xorm.Engine) *DeviceDAO {
	once.Do(func() {
		deviceDAO = &DeviceDAO{}
	})
	return deviceDAO
}

func (dao *DeviceDAO) InsertDevice(ctx context.Context,
	session *xorm.Session, device *DeviceDO) (*DeviceDO, error) {

	slog.InfoContext(ctx, "Inserting device", "device", device)
	rowsEffected, err := session.Insert(device)
	if err != nil {
		slog.ErrorContext(ctx, "insert device failed", "device",
			device, "error", err)
		return nil, err
	}
	if rowsEffected == 0 {
		errMsg := "no rows affected when inserting device"
		err = errors.New(errMsg)
		slog.ErrorContext(ctx, errMsg, "device_no", device.DeviceNo)
		return nil, err
	}
	slog.InfoContext(ctx, "Device inserted", "device", device,
		"rows_affected", rowsEffected)

	return device, nil
}

func (dao *DeviceDAO) QueryDeviceByDeviceNo(ctx context.Context,
	session *xorm.Session, deviceNo string) (*DeviceDO, error) {
	device := &DeviceDO{}
	has, err := session.Where("device_no = ?", deviceNo).Get(device)
	if err != nil {
		slog.ErrorContext(ctx, "query device by device_no failed", "device_no",
			deviceNo, "error", err)
		return nil, err
	}
	if !has {
		slog.InfoContext(ctx, "device not found", "device_no", deviceNo)
		return nil, nil
	}
	slog.InfoContext(ctx, "device found", "device", device)
	return device, nil
}

func (dao *DeviceDAO) QueryDeviceByDeviceId(ctx context.Context,
	session *xorm.Session, deviceId int64) (*DeviceDO, error) {
	device := &DeviceDO{}
	has, err := session.Where("id = ?", deviceId).Get(device)
	if err != nil {
		slog.ErrorContext(ctx, "query device by device_id failed", "device_id",
			deviceId, "error", err)
		return nil, err
	}
	if !has {
		slog.InfoContext(ctx, "device not found", "device_id", deviceId)
		return nil, nil
	}
	slog.InfoContext(ctx, "device found", "device", device)
	return device, nil
}
