package models

import (
	"time"

	"gorm.io/gorm"
)

type MealRecord struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	TotalPrice float64        `json:"total_price" gorm:"type:decimal(10,2);not null"`
	Thoughts   string         `json:"thoughts" gorm:"type:text"`
	ImageURL   string         `json:"image_url" gorm:"size:255"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	User   *User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Dishes []MealRecordDish `json:"dishes,omitempty" gorm:"foreignKey:MealRecordID"`
}

func (MealRecord) TableName() string {
	return "meal_records"
}
