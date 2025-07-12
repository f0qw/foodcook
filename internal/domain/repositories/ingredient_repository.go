package repositories

import (
	"context"

	"foodcook/internal/domain/models"
)

type IngredientRepository interface {
	Create(ctx context.Context, ingredient *models.Ingredient) error
	GetByID(ctx context.Context, id uint) (*models.Ingredient, error)
	Update(ctx context.Context, ingredient *models.Ingredient) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Ingredient, int64, error)
	IsUsedInDishes(ctx context.Context, ingredientID uint) (bool, error)
}
