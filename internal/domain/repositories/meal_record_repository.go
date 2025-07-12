package repositories

import (
	"context"

	"foodcook/internal/domain/models"
)

type MealRecordRepository interface {
	Create(ctx context.Context, mealRecord *models.MealRecord, dishIDs []uint) error
	GetByID(ctx context.Context, id uint) (*models.MealRecord, error)
	Update(ctx context.Context, mealRecord *models.MealRecord) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, userID uint, offset, limit int) ([]*models.MealRecord, int64, error)
	GetByUser(ctx context.Context, userID uint) ([]*models.MealRecord, error)
}
