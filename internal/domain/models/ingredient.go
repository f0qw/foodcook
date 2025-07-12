package models

import (
	"time"

	"gorm.io/gorm"
)

type Ingredient struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Price     float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	Unit      string         `json:"unit" gorm:"size:20;not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	DishIngredients []DishIngredient `json:"dish_ingredients,omitempty" gorm:"foreignKey:IngredientID"`
}

func (Ingredient) TableName() string {
	return "ingredients"
}
