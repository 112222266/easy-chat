package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用bcrypt算法对密码进行哈希处理
// 参数: password - 明文密码
// 返回: 哈希后的密码字符串和错误信息
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost 是推荐的加密成本(10)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 验证密码是否匹配哈希
// 参数: hashedPassword - 哈希后的密码, password - 待验证的明文密码
// 返回: 是否匹配
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
