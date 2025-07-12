package models

type DishIngredient struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	DishID       uint    `json:"dish_id" gorm:"not null"`
	IngredientID uint    `json:"ingredient_id" gorm:"not null"`
	Quantity     float64 `json:"quantity" gorm:"type:decimal(10,2);not null"`

	// 关联关系
	Dish       *Dish       `json:"dish,omitempty" gorm:"foreignKey:DishID"`
	Ingredient *Ingredient `json:"ingredient,omitempty" gorm:"foreignKey:IngredientID"`
}

func (DishIngredient) TableName() string {
	return "dish_ingredients"
}
