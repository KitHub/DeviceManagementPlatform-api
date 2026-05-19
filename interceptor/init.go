package interceptor

import (
	"DeviceManagementPlatform-api/servicecontext"
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	serviceName       = "gRPC-Jaeger-Demo"
	jaegerRPCEndpoint = "127.0.0.1:4317"
)

var tracer = otel.Tracer("grpc-example")

func init() {
	ctx := context.Background()
	tp, err := initTracer(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "init trace failed", slog.Any("error", err))
		os.Exit(1)
	}
	servicecontext.RegisterShutdownCallback(func(ctx context.Context) error {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
			return err
		}
		return nil
	})

}

func LogTraceInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	begin := time.Now()
	var traceId string
	if md, exist := metadata.FromIncomingContext(ctx); exist {
		if v, ok := md["traceid"]; ok {
			traceId = v[0]
		}
	}
	if traceId == "" {
		traceId = uuid.NewString()
	}
	res, err := handler(ctx, req)
	slog.InfoContext(ctx, "trace", slog.String("trace_id", traceId), slog.Int64("begin_time", begin.UnixMilli()), slog.Int64("duration", time.Since(begin).Milliseconds()))
	return res, err
}

func OtelTraceInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_, span := tracer.Start(ctx, "SayHello",
		trace.WithAttributes(
			attribute.StringSlice("client-id", md.Get("client-id")),
			attribute.StringSlice("user-id", md.Get("user-id")),
		),
	)
	defer span.End()
	return handler(ctx, req)
}

// initTracer 初始化 Tracer
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	tp, err := newJaegerTraceProvider(ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	return tp, nil
}

// newJaegerTraceProvider
func newJaegerTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(jaegerRPCEndpoint),
		otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
	if err != nil {
		return nil, err
	}
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(time.Second)),
	)
	return traceProvider, nil
}
