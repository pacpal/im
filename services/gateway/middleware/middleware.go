// Package middleware 提供 HTTP 层的 Gin 中间件：认证、CORS、日志与恢复等。
package middleware

import (
	"IM/pkg/auth"
	"IM/pkg/logger"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth 返回一个 Gin 中间件，用于解析 Authorization header 并验证 JWT，验证通过后将 user_id 注入到上下文中。
func Auth(jwtUtil *auth.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warnf("[HTTP Auth] missing authorization header method=%s path=%s ip=%s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warnf("[HTTP Auth] invalid authorization header method=%s path=%s ip=%s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		userID, err := jwtUtil.ParseToken(parts[1])
		if err != nil {
			logger.Warnf("[HTTP Auth] invalid token method=%s path=%s ip=%s error=%v", c.Request.Method, c.Request.URL.Path, c.ClientIP(), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("token", parts[1])
		c.Set("user_id", userID)
		c.Next()
	}
}

// CORS 返回一个允许跨域请求的中间件（仅用于简易开发场景）。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Logging 返回一个记录 HTTP 请求与响应状态的中间件。
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		if len(c.Errors) > 0 || status >= 400 {
			logger.Errorf("[HTTP] %s %s status=%d ip=%s errors=%s", method, path, status, c.ClientIP(), c.Errors.String())
		} else {
			logger.Infof("[HTTP] %s %s status=%d ip=%s", method, path, status, c.ClientIP())
		}
	}
}

// Recovery 返回一个 panic 恢复中间件，捕获 panic 并返回 500。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("[HTTP Recovery] panic=%v path=%s ip=%s\nstack:\n%s", r, c.Request.URL.Path, c.ClientIP(), debug.Stack())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}
