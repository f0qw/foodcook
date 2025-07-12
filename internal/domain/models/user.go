package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleUser = "user"
	RoleRoot = "root"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
	PasswordHash string         `json:"-" gorm:"size:255;not null"`
	Role         string         `json:"role" gorm:"size:20;default:'user';not null"` // user, root
	AvatarURL    string         `json:"avatar_url" gorm:"size:255"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	MealRecords []MealRecord `json:"meal_records,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
