package repositories

import (
	"context"

	"foodcook/internal/domain/models"
	"foodcook/internal/domain/repositories"

	"gorm.io/gorm"
)

type MySQLIngredientRepository struct {
	db *gorm.DB
}

func NewMySQLIngredientRepository(db *gorm.DB) repositories.IngredientRepository {
	return &MySQLIngredientRepository{db: db}
}

func (r *MySQLIngredientRepository) Create(ctx context.Context, ingredient *models.Ingredient) error {
	return r.db.WithContext(ctx).Create(ingredient).Error
}

func (r *MySQLIngredientRepository) GetByID(ctx context.Context, id uint) (*models.Ingredient, error) {
	var ingredient models.Ingredient
	err := r.db.WithContext(ctx).First(&ingredient, id).Error
	if err != nil {
		return nil, err
	}
	return &ingredient, nil
}

func (r *MySQLIngredientRepository) Update(ctx context.Context, ingredient *models.Ingredient) error {
	return r.db.WithContext(ctx).Save(ingredient).Error
}

func (r *MySQLIngredientRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Ingredient{}, id).Error
}

func (r *MySQLIngredientRepository) List(ctx context.Context, offset, limit int) ([]*models.Ingredient, int64, error) {
	var ingredients []*models.Ingredient
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Ingredient{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&ingredients).Error; err != nil {
		return nil, 0, err
	}

	return ingredients, total, nil
}

func (r *MySQLIngredientRepository) IsUsedInDishes(ctx context.Context, ingredientID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.DishIngredient{}).Where("ingredient_id = ?", ingredientID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
