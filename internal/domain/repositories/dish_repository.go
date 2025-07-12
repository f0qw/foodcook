package repositories

import (
	"context"

	"foodcook/internal/domain/models"
)

type DishRepository interface {
	Create(ctx context.Context, dish *models.Dish) error
	GetByID(ctx context.Context, id uint) (*models.Dish, error)
	Update(ctx context.Context, dish *models.Dish) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int, categoryID *uint) ([]*models.Dish, int64, error)
	GetByCategory(ctx context.Context, categoryID uint) ([]*models.Dish, error)
	Search(ctx context.Context, keyword string, offset, limit int) ([]*models.Dish, int64, error)
}
