package repositories

import (
	"context"
	"fmt"

	"foodcook/internal/domain/models"
	"foodcook/internal/domain/repositories"

	"gorm.io/gorm"
)

type MySQLCategoryRepository struct {
	db *gorm.DB
}

func NewMySQLCategoryRepository(db *gorm.DB) repositories.CategoryRepository {
	return &MySQLCategoryRepository{db: db}
}

func (r *MySQLCategoryRepository) Create(ctx context.Context, category *models.Category) error {
	result := r.db.WithContext(ctx).Create(category)
	if result.Error != nil {
		return fmt.Errorf("创建分类失败: %w", result.Error)
	}
	return nil
}

func (r *MySQLCategoryRepository) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category
	result := r.db.WithContext(ctx).First(&category, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("分类不存在")
		}
		return nil, fmt.Errorf("查询分类失败: %w", result.Error)
	}
	return &category, nil
}

func (r *MySQLCategoryRepository) Update(ctx context.Context, category *models.Category) error {
	result := r.db.WithContext(ctx).Save(category)
	if result.Error != nil {
		return fmt.Errorf("更新分类失败: %w", result.Error)
	}
	return nil
}

func (r *MySQLCategoryRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Category{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除分类失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("分类不存在")
	}
	return nil
}

func (r *MySQLCategoryRepository) List(ctx context.Context) ([]*models.Category, error) {
	var categories []*models.Category
	result := r.db.WithContext(ctx).Order("name ASC").Find(&categories)
	if result.Error != nil {
		return nil, fmt.Errorf("查询分类列表失败: %w", result.Error)
	}
	return categories, nil
}
