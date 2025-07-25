package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:50;unique;not null"`
	Email    string `gorm:"size:100;unique;not null"`
	Password string `gorm:"size:100;not null"`
	Age      int    `gorm:"default:18"`
	Birthday *time.Time
	Status   bool `gorm:"default:true"`
}

// TableName 设置表名
func (User) TableName() string {
	return "user"
}
