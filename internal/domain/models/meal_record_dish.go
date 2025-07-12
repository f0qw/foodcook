package models

type MealRecordDish struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	MealRecordID uint `json:"meal_record_id" gorm:"not null"`
	DishID       uint `json:"dish_id" gorm:"not null"`
	Quantity     int  `json:"quantity" gorm:"default:1"`

	// 关联关系
	MealRecord *MealRecord `json:"meal_record,omitempty" gorm:"foreignKey:MealRecordID"`
	Dish       *Dish       `json:"dish,omitempty" gorm:"foreignKey:DishID"`
}

func (MealRecordDish) TableName() string {
	return "meal_record_dishes"
}
