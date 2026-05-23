// Package auth 提供 JWT 生成与解析的简单封装，便于在服务间进行身份认证。
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTUtil 包装了用于生成与解析 JWT 的密钥与过期配置。
type JWTUtil struct {
	secretKey  []byte
	expiration time.Duration
}

// NewJWTUtil 使用指定的 secretKey 与过期时长创建 JWTUtil。
func NewJWTUtil(secretKey string, expiration time.Duration) *JWTUtil {
	return &JWTUtil{
		secretKey:  []byte(secretKey),
		expiration: expiration,
	}
}

// Claims 定义自定义的 JWT Claims，包含用户 ID 与用户名。
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 为指定用户生成签名的 token。
func (j *JWTUtil) GenerateToken(userID, username string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ParseToken 使用 JWTUtil 的密钥解析 token 并返回 userID（如验证失败返回错误）。
func (j *JWTUtil) ParseToken(tokenStr string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", jwt.ErrSignatureInvalid
	}
	return claims.UserID, nil
}

// GetSecret 返回底层密钥字节切片。
func (j *JWTUtil) GetSecret() []byte {
	return j.secretKey
}

// GetExpiration 返回配置的过期时长。
func (j *JWTUtil) GetExpiration() time.Duration {
	return j.expiration
}

// GenerateToken 为给定 secret 和过期时间生成 JWT（独立函数，便于无需构造 JWTUtil 时使用）。
func GenerateToken(userID, username string, secret []byte, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken 解析并返回 token 中的 Claims（独立函数）。
func ParseToken(tokenStr string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
