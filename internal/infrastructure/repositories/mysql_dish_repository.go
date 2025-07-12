package repositories

import (
	"context"
	"fmt"

	"foodcook/internal/domain/models"
	"foodcook/internal/domain/repositories"

	"gorm.io/gorm"
)

type MySQLDishRepository struct {
	db *gorm.DB
}

func NewMySQLDishRepository(db *gorm.DB) repositories.DishRepository {
	return &MySQLDishRepository{db: db}
}

func (r *MySQLDishRepository) List(ctx context.Context, offset, limit int, categoryID *uint) ([]*models.Dish, int64, error) {
	var dishes []*models.Dish
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Dish{}).Preload("Category").Preload("Ingredients.Ingredient")

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询菜品总数失败: %w", err)
	}

	// 查询菜品列表
	result := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&dishes)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("查询菜品列表失败: %w", result.Error)
	}

	return dishes, total, nil
}

func (r *MySQLDishRepository) GetByID(ctx context.Context, id uint) (*models.Dish, error) {
	var dish models.Dish
	result := r.db.WithContext(ctx).Preload("Category").Preload("Ingredients.Ingredient").First(&dish, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("菜品不存在")
		}
		return nil, fmt.Errorf("查询菜品失败: %w", result.Error)
	}
	return &dish, nil
}

func (r *MySQLDishRepository) Create(ctx context.Context, dish *models.Dish) error {
	result := r.db.WithContext(ctx).Create(dish)
	if result.Error != nil {
		return fmt.Errorf("创建菜品失败: %w", result.Error)
	}
	return nil
}

func (r *MySQLDishRepository) CreateWithIngredients(ctx context.Context, dish *models.Dish, ingredients []repositories.DishIngredientRequest) error {
	// 开始事务
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建菜品
	if err := tx.Create(dish).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建菜品失败: %w", err)
	}

	// 创建食材关联
	for _, ingredient := range ingredients {
		dishIngredient := &models.DishIngredient{
			DishID:       dish.ID,
			IngredientID: ingredient.IngredientID,
			Quantity:     ingredient.Quantity,
		}
		if err := tx.Create(dishIngredient).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建食材关联失败: %w", err)
		}
	}

	// 提交事务
	return tx.Commit().Error
}

func (r *MySQLDishRepository) Update(ctx context.Context, dish *models.Dish) error {
	result := r.db.WithContext(ctx).Save(dish)
	if result.Error != nil {
		return fmt.Errorf("更新菜品失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("菜品不存在")
	}
	return nil
}

func (r *MySQLDishRepository) UpdateWithIngredients(ctx context.Context, dish *models.Dish, ingredients []repositories.DishIngredientRequest) error {
	// 开始事务
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新菜品
	if err := tx.Save(dish).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新菜品失败: %w", err)
	}

	// 删除旧的食材关联
	if err := tx.Where("dish_id = ?", dish.ID).Delete(&models.DishIngredient{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除旧食材关联失败: %w", err)
	}

	// 创建新的食材关联
	for _, ingredient := range ingredients {
		dishIngredient := &models.DishIngredient{
			DishID:       dish.ID,
			IngredientID: ingredient.IngredientID,
			Quantity:     ingredient.Quantity,
		}
		if err := tx.Create(dishIngredient).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建食材关联失败: %w", err)
		}
	}

	// 提交事务
	return tx.Commit().Error
}

func (r *MySQLDishRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Dish{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除菜品失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("菜品不存在")
	}
	return nil
}

func (r *MySQLDishRepository) GetByCategory(ctx context.Context, categoryID uint) ([]*models.Dish, error) {
	var dishes []*models.Dish
	result := r.db.WithContext(ctx).Preload("Category").Where("category_id = ?", categoryID).Order("created_at DESC").Find(&dishes)
	if result.Error != nil {
		return nil, fmt.Errorf("查询分类菜品失败: %w", result.Error)
	}
	return dishes, nil
}

func (r *MySQLDishRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*models.Dish, int64, error) {
	var dishes []*models.Dish
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Dish{}).Preload("Category").
		Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询搜索结果总数失败: %w", err)
	}

	// 查询搜索结果
	result := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&dishes)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("搜索菜品失败: %w", result.Error)
	}

	return dishes, total, nil
}

func (r *MySQLDishRepository) IsUsedInMealRecords(ctx context.Context, dishID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.MealRecordDish{}).
		Joins("JOIN meal_records ON meal_records.id = meal_record_dishes.meal_record_id").
		Where("meal_record_dishes.dish_id = ? AND meal_records.deleted_at IS NULL", dishID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查菜品使用情况失败: %w", err)
	}
	return count > 0, nil
}
