package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义 JWT claims 结构体
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
// 参数:
//
//	secret - JWT 签名密钥
//	expire - 过期时间(秒)
//	userID - 用户ID(uint类型)
//
// 返回:
//
//	token 字符串和错误信息
func GenerateToken(secret string, expire int, userID uint) (string, error) {
	// 设置 token 过期时间
	expirationTime := time.Now().Add(time.Duration(expire) * time.Second)

	// 创建 claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析和验证 JWT token
// 参数:
//
//	tokenString - token 字符串
//	secret - JWT 签名密钥
//
// 返回:
//
//	claims 和错误信息
func ParseToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
