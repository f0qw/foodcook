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

	query := r.db.WithContext(ctx).Model(&models.Dish{}).Preload("Category")

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
	result := r.db.WithContext(ctx).Preload("Category").First(&dish, id)
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
