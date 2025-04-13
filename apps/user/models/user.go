package models

import "gorm.io/gorm"

// User 用户模型
type User struct {
	gorm.Model
	Avatar   string `gorm:"type:varchar(191);not null;default:'';comment:用户头像URL"`
	Nickname string `gorm:"type:varchar(24);not null;comment:用户昵称"`
	Phone    string `gorm:"type:varchar(20);not null;uniqueIndex:idx_phone;comment:手机号码"`
	Password string `gorm:"type:varchar(191);not null;comment:密码(加密存储)"`
	Status   int8   `gorm:"type:tinyint;not null;default:1;comment:用户状态(0-禁用,1-正常,2-锁定)"`
	Sex      int8   `gorm:"type:tinyint;not null;default:0;comment:性别(0-未知,1-男,2-女)"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}
