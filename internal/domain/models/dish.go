package models

import (
	"time"

	"gorm.io/gorm"
)

type Dish struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Description string         `json:"description" gorm:"type:text"`
	ImageURL    string         `json:"image_url" gorm:"size:255"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	CookingLink string         `json:"cooking_link" gorm:"size:255"`
	CategoryID  *uint          `json:"category_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Category    *Category        `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Ingredients []DishIngredient `json:"ingredients,omitempty" gorm:"foreignKey:DishID"`
	MealRecords []MealRecordDish `json:"meal_records,omitempty" gorm:"foreignKey:DishID"`
}

func (Dish) TableName() string {
	return "dishes"
}
