package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret")

func DBCreateID(name string) string {
	return "DBkey"
}

// 注册

func GenerateToken(name, ID string) (string, error) {
	lastTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"name": name,
		"sub":  ID,
		"exp":  lastTime.Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

// 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "'missing token'"})
			c.Abort()
			return
		}
		parts := strings.Split(header, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			c.Abort()
			return
		}
		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization"})
			c.Abort()
			return
		}

		c.Set("userID", claims["sub"])
		c.Set("userName", claims["name"])
		c.Next()
	}
}

// token 解析
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("jwt internal wrong")
			}
			return jwtSecret, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithNotBeforeRequired(),
		// jwt.WithIssuer("servername"),
	)
	if err != nil {
		// error type make choice.
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token invalid")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, nil
}
