package dao

import "time"

type DeviceDO struct {
	Id           int64     `xorm:"id,pk,autoincr"`
	DeviceNo     string    `xorm:"device_no"`
	RegisterTime time.Time `xorm:"register_time"`
	CreateTime   time.Time `xorm:"create_time"`
	UpdateTime   time.Time `xorm:"update_time"`
}

func (d *DeviceDO) TableName() string {
	return "device"
}
