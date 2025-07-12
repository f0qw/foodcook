package handlers

import (
	"net/http"
	"strconv"

	"foodcook/internal/domain/models"
	"foodcook/internal/domain/repositories"

	"github.com/gin-gonic/gin"
)

type DishHandler struct {
	dishRepo repositories.DishRepository
}

func NewDishHandler(dishRepo repositories.DishRepository) *DishHandler {
	return &DishHandler{
		dishRepo: dishRepo,
	}
}

type CreateDishRequest struct {
	Name        string                  `json:"name" binding:"required"`
	Description string                  `json:"description"`
	ImageURL    string                  `json:"image_url"`
	Price       float64                 `json:"price" binding:"required,min=0"`
	CookingLink string                  `json:"cooking_link"`
	CategoryID  *uint                   `json:"category_id"`
	Ingredients []DishIngredientRequest `json:"ingredients"`
}

type UpdateDishRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ImageURL    string                  `json:"image_url"`
	Price       float64                 `json:"price" binding:"min=0"`
	CookingLink string                  `json:"cooking_link"`
	CategoryID  *uint                   `json:"category_id"`
	Ingredients []DishIngredientRequest `json:"ingredients"`
}

type DishIngredientRequest struct {
	IngredientID uint    `json:"ingredient_id" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required,min=0"`
}

func (h *DishHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	categoryIDStr := c.Query("category_id")

	var categoryID *uint
	if categoryIDStr != "" {
		if id, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			catID := uint(id)
			categoryID = &catID
		}
	}

	dishes, total, err := h.dishRepo.List(c.Request.Context(), offset, limit, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜品列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   dishes,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

func (h *DishHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜品ID"})
		return
	}

	dish, err := h.dishRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜品不存在"})
		return
	}

	c.JSON(http.StatusOK, dish)
}

func (h *DishHandler) Create(c *gin.Context) {
	var req CreateDishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dish := &models.Dish{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		CookingLink: req.CookingLink,
		CategoryID:  req.CategoryID,
	}

	// 转换食材请求为repository类型
	var ingredients []repositories.DishIngredientRequest
	for _, ing := range req.Ingredients {
		ingredients = append(ingredients, repositories.DishIngredientRequest{
			IngredientID: ing.IngredientID,
			Quantity:     ing.Quantity,
		})
	}

	// 创建菜品和食材关联
	if err := h.dishRepo.CreateWithIngredients(c.Request.Context(), dish, ingredients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建菜品失败"})
		return
	}

	c.JSON(http.StatusCreated, dish)
}

func (h *DishHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜品ID"})
		return
	}

	var req UpdateDishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dish, err := h.dishRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜品不存在"})
		return
	}

	// 更新字段
	if req.Name != "" {
		dish.Name = req.Name
	}
	if req.Description != "" {
		dish.Description = req.Description
	}
	if req.ImageURL != "" {
		dish.ImageURL = req.ImageURL
	}
	if req.Price > 0 {
		dish.Price = req.Price
	}
	if req.CookingLink != "" {
		dish.CookingLink = req.CookingLink
	}
	if req.CategoryID != nil {
		dish.CategoryID = req.CategoryID
	}

	// 转换食材请求为repository类型
	var ingredients []repositories.DishIngredientRequest
	for _, ing := range req.Ingredients {
		ingredients = append(ingredients, repositories.DishIngredientRequest{
			IngredientID: ing.IngredientID,
			Quantity:     ing.Quantity,
		})
	}

	// 更新菜品和食材关联
	if err := h.dishRepo.UpdateWithIngredients(c.Request.Context(), dish, ingredients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新菜品失败"})
		return
	}

	c.JSON(http.StatusOK, dish)
}

func (h *DishHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜品ID"})
		return
	}

	if err := h.dishRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除菜品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "菜品删除成功"})
}

func (h *DishHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	dishes, total, err := h.dishRepo.Search(c.Request.Context(), keyword, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索菜品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   dishes,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}
