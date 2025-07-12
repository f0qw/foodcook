package handlers

import (
	"net/http"
	"strconv"

	"foodcook/internal/domain/models"
	"foodcook/internal/domain/repositories"

	"github.com/gin-gonic/gin"
)

type MealRecordHandler struct {
	mealRecordRepo repositories.MealRecordRepository
	dishRepo       repositories.DishRepository
}

func NewMealRecordHandler(mealRecordRepo repositories.MealRecordRepository, dishRepo repositories.DishRepository) *MealRecordHandler {
	return &MealRecordHandler{
		mealRecordRepo: mealRecordRepo,
		dishRepo:       dishRepo,
	}
}

type CreateMealRecordRequest struct {
	DishIDs  []uint `json:"dish_ids" binding:"required"`
	Thoughts string `json:"thoughts"`
	ImageURL string `json:"image_url"`
}

func (h *MealRecordHandler) List(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	mealRecords, total, err := h.mealRecordRepo.List(c.Request.Context(), userID.(uint), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用餐记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   mealRecords,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

func (h *MealRecordHandler) GetByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用餐记录ID"})
		return
	}

	mealRecord, err := h.mealRecordRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用餐记录不存在"})
		return
	}

	// 检查权限
	if mealRecord.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此用餐记录"})
		return
	}

	c.JSON(http.StatusOK, mealRecord)
}

func (h *MealRecordHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req CreateMealRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 计算总价格
	var totalPrice float64
	for _, dishID := range req.DishIDs {
		dish, err := h.dishRepo.GetByID(c.Request.Context(), dishID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "菜品不存在"})
			return
		}
		totalPrice += dish.Price
	}

	mealRecord := &models.MealRecord{
		UserID:     userID.(uint),
		TotalPrice: totalPrice,
		Thoughts:   req.Thoughts,
		ImageURL:   req.ImageURL,
	}

	if err := h.mealRecordRepo.Create(c.Request.Context(), mealRecord, req.DishIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用餐记录失败"})
		return
	}

	c.JSON(http.StatusCreated, mealRecord)
}

func (h *MealRecordHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用餐记录ID"})
		return
	}

	mealRecord, err := h.mealRecordRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用餐记录不存在"})
		return
	}

	// 检查权限
	if mealRecord.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此用餐记录"})
		return
	}

	if err := h.mealRecordRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用餐记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用餐记录删除成功"})
}
