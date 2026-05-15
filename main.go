package main

import (
	"DeviceManagementPlatform-api/config"
	servicecontext "DeviceManagementPlatform-api/service_context"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/KitHub/protocols/devicemanagementplatformapi"
	"google.golang.org/grpc"
)

type ServerArgs struct {
	ConfigFile string
}

type wrappedServer struct {
	*grpc.Server
	listerners net.Listener
}

// Close implements [io.Closer].
func (w *wrappedServer) Close() error {
	w.GracefulStop()
	return nil
}

func main() {
	ctx := context.Background()
	args := parepareArgs(ctx)

	// init config
	slog.InfoContext(ctx, "init config",
		slog.String("config_file", args.ConfigFile))
	configEntity, err := config.LoadConfig(ctx, args.ConfigFile)
	if err != nil {
		slog.ErrorContext(ctx, "init config failed",
			slog.String("error", err.Error()))
		panic(err)
	}
	slog.InfoContext(ctx, "init config done")

	// init service context
	serviceContext, err := servicecontext.InitServiceContext(
		ctx, &configEntity)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init service context",
			slog.String("error", err.Error()))
		panic(err)
	}

	// init server
	servers, err := initServer(ctx, &configEntity, serviceContext)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init server",
			slog.String("error", err.Error()))
		panic(err)
	}

	// start server
	for _, srv := range servers {
		go srv.Serve(srv.listerners)
	}

	closers := make([]io.Closer, 0)
	for _, srv := range servers {
		closers = append(closers, srv)
	}
	addShutdownHook(ctx, closers)
}

func parepareArgs(ctx context.Context) ServerArgs {
	configFile := flag.String("server_config", "", "config file for server")
	flag.Parse()
	result := ServerArgs{
		ConfigFile: *configFile,
	}
	slog.InfoContext(ctx, "parse flags done")
	return result
}

func initServer(ctx context.Context, serviceConfig *config.ConfigEntity,
	serviceContext *servicecontext.ServiceContext) (
	servers []*wrappedServer, err error) {

	slog.InfoContext(ctx, "init servers")
	for _, serverConfig := range serviceConfig.Servers {
		hostAndPort := fmt.Sprintf("%s:%d",
			serverConfig.Host, serverConfig.Port)
		slog.InfoContext(ctx, "init server",
			slog.String("host_and_port", hostAndPort),
			slog.String("type", serverConfig.Type))
		listener, err := net.Listen("tcp", hostAndPort)
		if err != nil {
			slog.ErrorContext(ctx, "failed to listen",
				slog.String("error", err.Error()))
			return nil, err
		}

		// create a new gRPC server
		server := grpc.NewServer()
		// bind the service implementation to the gRPC server
		devicemanagementplatformapi.RegisterDeviceManagementPlatformAPIServer(
			server, serviceContext.ApiService)

		wrappedSrv := &wrappedServer{
			Server:     server,
			listerners: listener,
		}
		servers = append(servers, wrappedSrv)
	}
	slog.InfoContext(ctx, "init servers done")

	return servers, nil
}
func addShutdownHook(ctx context.Context, closers []io.Closer) {
	slog.InfoContext(ctx, "listening signals...")
	c := make(chan os.Signal, 1)
	signal.Notify(
		c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	<-c
	slog.InfoContext(ctx, "graceful shutdown...")

	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			slog.ErrorContext(ctx, "failed to stop closer", slog.Any("error", err))
		}
	}

	slog.InfoContext(ctx, "completed graceful shutdown")
}
