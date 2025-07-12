package handlers

import (
	"net/http"
	"strconv"

	"foodcook/internal/domain/models"
	"foodcook/internal/domain/repositories"

	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	ingredientRepo repositories.IngredientRepository
}

func NewIngredientHandler(ingredientRepo repositories.IngredientRepository) *IngredientHandler {
	return &IngredientHandler{
		ingredientRepo: ingredientRepo,
	}
}

type CreateIngredientRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,min=0"`
	Unit  string  `json:"unit" binding:"required"`
}

type UpdateIngredientRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price" binding:"min=0"`
	Unit  string  `json:"unit"`
}

func (h *IngredientHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	ingredients, total, err := h.ingredientRepo.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取食材列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   ingredients,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

func (h *IngredientHandler) Create(c *gin.Context) {
	var req CreateIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredient := &models.Ingredient{
		Name:  req.Name,
		Price: req.Price,
		Unit:  req.Unit,
	}

	if err := h.ingredientRepo.Create(c.Request.Context(), ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建食材失败"})
		return
	}

	c.JSON(http.StatusCreated, ingredient)
}

func (h *IngredientHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的食材ID"})
		return
	}

	var req UpdateIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ingredient, err := h.ingredientRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "食材不存在"})
		return
	}

	// 更新字段
	if req.Name != "" {
		ingredient.Name = req.Name
	}
	if req.Price > 0 {
		ingredient.Price = req.Price
	}
	if req.Unit != "" {
		ingredient.Unit = req.Unit
	}

	if err := h.ingredientRepo.Update(c.Request.Context(), ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新食材失败"})
		return
	}

	c.JSON(http.StatusOK, ingredient)
}

func (h *IngredientHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的食材ID"})
		return
	}

	// 检查食材是否被菜品使用
	isUsed, err := h.ingredientRepo.IsUsedInDishes(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查食材使用情况失败"})
		return
	}

	if isUsed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该食材已被菜品使用，无法删除"})
		return
	}

	if err := h.ingredientRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除食材失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "食材删除成功"})
}
