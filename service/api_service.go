package service

import (
	"DeviceManagementPlatform-api/logic"
	"context"
	"sync"

	"github.com/KitHub/protocols/devicemanagementplatformapi"
)

type ApiService struct {
	devicemanagementplatformapi.UnimplementedDeviceManagementPlatformAPIServer
	deviceLogic *logic.DeviceLogic
}

var apiSrv *ApiService
var once sync.Once

func NewApiService(ctx context.Context,
	deviceLogic *logic.DeviceLogic) *ApiService {
	once.Do(func() {
		apiSrv = &ApiService{
			deviceLogic: deviceLogic,
		}
	})
	return apiSrv
}

// QueryDeviceById implements [devicemanagementplatformapi.DeviceManagementPlatformAPIServer].
func (a *ApiService) QueryDeviceById(context.Context, *devicemanagementplatformapi.QueryDeviceByIdRequest) (*devicemanagementplatformapi.QueryDeviceByIdResponse, error) {
	panic("unimplemented")
}

// QueryDeviceByNo implements [devicemanagementplatformapi.DeviceManagementPlatformAPIServer].
func (a *ApiService) QueryDeviceByNo(context.Context, *devicemanagementplatformapi.QueryDeviceByNoRequest) (*devicemanagementplatformapi.QueryDeviceByNoResponse, error) {
	panic("unimplemented")
}

// RegisterDevice implements [devicemanagementplatformapi.DeviceManagementPlatformAPIServer].
func (a *ApiService) RegisterDevice(context.Context, *devicemanagementplatformapi.RegisterDeviceRequest) (*devicemanagementplatformapi.RegisterDeviceResponse, error) {
	panic("unimplemented")
}
