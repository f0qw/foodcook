package routes

import (
	"foodcook/internal/app/handlers"
	"foodcook/internal/app/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	authHandler *handlers.AuthHandler,
	dishHandler *handlers.DishHandler,
	ingredientHandler *handlers.IngredientHandler,
	mealRecordHandler *handlers.MealRecordHandler,
	categoryHandler *handlers.CategoryHandler,
) *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggingMiddleware())

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
		}

		// 菜品路由 - 只有 root 用户可以管理
		dishes := api.Group("/dishes")
		{
			dishes.GET("", dishHandler.List)          // 所有用户都可以查看菜品列表
			dishes.GET("/:id", dishHandler.GetByID)   // 所有用户都可以查看菜品详情
			dishes.GET("/search", dishHandler.Search) // 所有用户都可以搜索菜品
			// 以下操作需要 root 权限
			dishes.POST("", middleware.AuthMiddleware(), middleware.RootMiddleware(), dishHandler.Create)
			dishes.PUT("/:id", middleware.AuthMiddleware(), middleware.RootMiddleware(), dishHandler.Update)
			dishes.DELETE("/:id", middleware.AuthMiddleware(), middleware.RootMiddleware(), dishHandler.Delete)
		}

		// 食材路由 - 只有 root 用户可以管理
		ingredients := api.Group("/ingredients")
		{
			ingredients.GET("", ingredientHandler.List) // 所有用户都可以查看食材列表
			// 以下操作需要 root 权限
			ingredients.POST("", middleware.AuthMiddleware(), middleware.RootMiddleware(), ingredientHandler.Create)
			ingredients.PUT("/:id", middleware.AuthMiddleware(), middleware.RootMiddleware(), ingredientHandler.Update)
			ingredients.DELETE("/:id", middleware.AuthMiddleware(), middleware.RootMiddleware(), ingredientHandler.Delete)
		}

		// 分类路由 - 只有 root 用户可以管理
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.List) // 所有用户都可以查看分类列表
			// 以下操作需要 root 权限
			categories.POST("", middleware.AuthMiddleware(), middleware.RootMiddleware(), categoryHandler.Create)
			categories.PUT("/:id", middleware.AuthMiddleware(), middleware.RootMiddleware(), categoryHandler.Update)
			categories.DELETE("/:id", middleware.AuthMiddleware(), middleware.RootMiddleware(), categoryHandler.Delete)
		}

		// 用餐记录路由
		mealRecords := api.Group("/meal-records")
		{
			mealRecords.GET("", middleware.AuthMiddleware(), mealRecordHandler.List)
			mealRecords.POST("", middleware.AuthMiddleware(), mealRecordHandler.Create)
			mealRecords.GET("/:id", middleware.AuthMiddleware(), mealRecordHandler.GetByID)
			mealRecords.PUT("/:id", middleware.AuthMiddleware(), mealRecordHandler.Update)
			mealRecords.DELETE("/:id", middleware.AuthMiddleware(), mealRecordHandler.Delete)
		}
	}

	return r
}
