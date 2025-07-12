package repositories

import (
	"context"

	"foodcook/internal/domain/models"
	"foodcook/internal/domain/repositories"

	"gorm.io/gorm"
)

type MySQLMealRecordRepository struct {
	db *gorm.DB
}

func NewMySQLMealRecordRepository(db *gorm.DB) repositories.MealRecordRepository {
	return &MySQLMealRecordRepository{db: db}
}

func (r *MySQLMealRecordRepository) Create(ctx context.Context, mealRecord *models.MealRecord, dishIDs []uint) error {
	// 开始事务
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用餐记录
	if err := tx.Create(mealRecord).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建菜品关联
	for _, dishID := range dishIDs {
		mealRecordDish := &models.MealRecordDish{
			MealRecordID: mealRecord.ID,
			DishID:       dishID,
			Quantity:     1, // 默认数量为1
		}
		if err := tx.Create(mealRecordDish).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	return tx.Commit().Error
}

func (r *MySQLMealRecordRepository) GetByID(ctx context.Context, id uint) (*models.MealRecord, error) {
	var mealRecord models.MealRecord
	err := r.db.WithContext(ctx).Preload("Dishes.Dish").First(&mealRecord, id).Error
	if err != nil {
		return nil, err
	}
	return &mealRecord, nil
}

func (r *MySQLMealRecordRepository) Update(ctx context.Context, mealRecord *models.MealRecord) error {
	return r.db.WithContext(ctx).Save(mealRecord).Error
}

func (r *MySQLMealRecordRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.MealRecord{}, id).Error
}

func (r *MySQLMealRecordRepository) List(ctx context.Context, userID uint, offset, limit int) ([]*models.MealRecord, int64, error) {
	var mealRecords []*models.MealRecord
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.MealRecord{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，包含菜品信息
	if err := r.db.WithContext(ctx).Preload("Dishes.Dish").Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&mealRecords).Error; err != nil {
		return nil, 0, err
	}

	return mealRecords, total, nil
}

func (r *MySQLMealRecordRepository) GetByUser(ctx context.Context, userID uint) ([]*models.MealRecord, error) {
	var mealRecords []*models.MealRecord
	err := r.db.WithContext(ctx).Preload("Dishes.Dish").Where("user_id = ?", userID).Order("created_at DESC").Find(&mealRecords).Error
	if err != nil {
		return nil, err
	}
	return mealRecords, nil
}
