package interceptor

import (
	"context"
	"github.com/stawwkom/auth_service/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"time"
)

func LogInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// Вытаскиваем peer info (IP клиента)
	var clientIP string
	if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}

	// Выполняем обработку запроса
	res, err := handler(ctx, req)

	duration := time.Since(start)

	fields := []zap.Field{
		zap.String("method", info.FullMethod),
		zap.String("client_ip", clientIP),
		zap.Duration("duration", duration),
		zap.Any("request", req),
		zap.Any("response", res),
	}

	// Пример: логируем user_id, если он в ctx (например, после валидации JWT)
	if userID, ok := ctx.Value("user_id").(string); ok {
		fields = append(fields, zap.String("user_id", userID))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		logger.Error("gRPC request failed", fields...)
	} else {
		logger.Info("gRPC request handled", fields...)
	}

	return res, err
}
