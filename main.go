package main

import (
	"DeviceManagementPlatform-api/config"
	servicecontext "DeviceManagementPlatform-api/service_context"
	"context"
	"flag"
)

type ServerArgs struct {
	ConfigFile string
}

func main() {
	ctx := context.Background()
	args := parepareArgs(ctx)

	configEntity, err := config.LoadConfig(ctx, args.ConfigFile)
	if err != nil {
		panic(err)
	}

	_, err = servicecontext.InitServiceContext(ctx, configEntity)
	if err != nil {
		panic(err)
	}
}

func parepareArgs(ctx context.Context) ServerArgs {
	configFile := flag.String("server_config", "", "config file for server")
	flag.Parse()
	return ServerArgs{
		ConfigFile: *configFile,
	}
}
