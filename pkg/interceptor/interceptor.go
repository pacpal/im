// Package interceptor 提供一组 gRPC 拦截器：日志记录、panic 恢复与认证注入等。
package interceptor

import (
	"context"
	"runtime/debug"
	"strings"
	"time"

	"IM/pkg/auth"
	"IM/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

// AuthUnaryInterceptor 认证拦截器，从 metadata 提取 token 并校验后注入 ctx
// skipMethods 为不需要认证的方法列表
func AuthUnaryInterceptor(jwtSecret []byte, skipMethods []string) grpc.UnaryServerInterceptor {
	skipMap := make(map[string]bool, len(skipMethods))
	for _, m := range skipMethods {
		skipMap[m] = true
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 跳过不需要认证的方法
		if skipMap[info.FullMethod] {
			return handler(ctx, req)
		}

		// 从 metadata 提取 token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
		}

		token := authHeader[0]
		// 支持 "Bearer <token>" 格式
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// 解析 token
		claims, err := auth.ParseToken(token, jwtSecret)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// 将用户信息注入 context
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)

		return handler(ctx, req)
	}
}
