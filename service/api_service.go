package service

import (
	"DeviceManagementPlatform-api/logic"
	"context"
	"log/slog"
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
func (a *ApiService) QueryDeviceById(ctx context.Context, req *devicemanagementplatformapi.QueryDeviceByIdRequest) (rsp *devicemanagementplatformapi.QueryDeviceByIdResponse, err error) {
	slog.InfoContext(ctx, "Querying device by ID", slog.Any("deviceId", req.GetDeviceId()))

	deviceInfo, err := a.deviceLogic.GetDeviceById(ctx, req.GetDeviceId())
	if err != nil {
		slog.ErrorContext(ctx, "Failed to query device by ID", slog.Any("error", err))
		return nil, err
	}

	rsp = &devicemanagementplatformapi.QueryDeviceByIdResponse{
		Data: &devicemanagementplatformapi.QueryDeviceByIdResponseData{
			DeviceInfo: &devicemanagementplatformapi.DeviceInfo{
				DeviceId:     deviceInfo.Id,
				DeviceNo:     deviceInfo.DeviceNo,
				RegisterTime: deviceInfo.RegisterTime.UnixMilli(),
			},
		},
	}
	slog.InfoContext(ctx, "Successfully queried device by ID", slog.Any("device", deviceInfo))
	return rsp, nil
}

// QueryDeviceByNo implements [devicemanagementplatformapi.DeviceManagementPlatformAPIServer].
func (a *ApiService) QueryDeviceByNo(ctx context.Context, req *devicemanagementplatformapi.QueryDeviceByNoRequest) (rsp *devicemanagementplatformapi.QueryDeviceByNoResponse, err error) {
	slog.InfoContext(ctx, "Querying device by NO", slog.Any("deviceNo", req.GetDeviceNo()))

	deviceInfo, err := a.deviceLogic.GetDeviceByNo(ctx, req.GetDeviceNo())
	if err != nil {
		slog.ErrorContext(ctx, "Failed to query device by NO", slog.Any("error", err))
		return nil, err
	}

	rsp = &devicemanagementplatformapi.QueryDeviceByNoResponse{
		Data: &devicemanagementplatformapi.QueryDeviceByNoResponseData{
			DeviceInfo: &devicemanagementplatformapi.DeviceInfo{
				DeviceId:     deviceInfo.Id,
				DeviceNo:     deviceInfo.DeviceNo,
				RegisterTime: deviceInfo.RegisterTime.UnixMilli(),
			},
		},
	}
	slog.InfoContext(ctx, "Successfully queried device by NO", slog.Any("device", deviceInfo))
	return rsp, nil
}

// RegisterDevice implements [devicemanagementplatformapi.DeviceManagementPlatformAPIServer].
func (a *ApiService) RegisterDevice(ctx context.Context, req *devicemanagementplatformapi.RegisterDeviceRequest) (rsp *devicemanagementplatformapi.RegisterDeviceResponse, err error) {
	slog.InfoContext(ctx, "Registering device", slog.Any("deviceNo", req.GetDeviceNo()))

	deviceInfo, err := a.deviceLogic.RegisterDevice(ctx, req.GetDeviceNo())
	if err != nil {
		slog.ErrorContext(ctx, "Failed to register device", slog.Any("error", err))
		return nil, err
	}

	rsp = &devicemanagementplatformapi.RegisterDeviceResponse{
		Data: &devicemanagementplatformapi.RegisterDeviceResponseData{
			DeviceInfo: &devicemanagementplatformapi.DeviceInfo{
				DeviceId:     deviceInfo.Id,
				DeviceNo:     deviceInfo.DeviceNo,
				RegisterTime: deviceInfo.RegisterTime.UnixMilli(),
			},
		},
	}
	slog.InfoContext(ctx, "Successfully registered device", slog.Any("device", deviceInfo))
	return rsp, nil
}
