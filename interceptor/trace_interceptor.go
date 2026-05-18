package interceptor

import (
	"context"
)

func TraceInterceptor(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error) {
	// 这里可以添加 OpenTelemetry 的相关代码来创建和管理 TracerProvider
	// 例如，您可以在这里初始化一个全局的 TracerProvider，并在每个请求中使用它来创建 Span
}
