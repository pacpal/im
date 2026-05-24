// Package interceptor 提供一组 gRPC 拦截器：日志记录、panic 恢复与认证注入等。
package interceptor

import (
	"IM/pkg/logger"
	"context"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingUnaryInterceptor 返回一个 UnaryServerInterceptor，用于记录方法名和耗时。
func LoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			logger.Errorw("gRPC request completed",
				"component", "grpc", "method", info.FullMethod, "duration", duration, "error", err)
		} else {
			logger.Infow("gRPC request completed",
				"component", "grpc", "method", info.FullMethod, "duration", duration)
		}

		return resp, err
	}
}

// RecoveryUnaryInterceptor panic 恢复拦截器
func RecoveryUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorw("gRPC panic recovered",
					"component", "grpc_recovery", "method", info.FullMethod, "panic", r, "stack", string(debug.Stack()))
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

type contextKey string

const (
	// UserIDKey 存储在 context 中的 user ID key
	UserIDKey contextKey = "user_id"
	// UsernameKey 存储在 context 中的 username key
	UsernameKey contextKey = "username"
)
