package routes

import (
	"foodcook/internal/app/handlers"
	"foodcook/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	authHandler *handlers.AuthHandler,
	dishHandler *handlers.DishHandler,
	ingredientHandler *handlers.IngredientHandler,
	mealRecordHandler *handlers.MealRecordHandler,
) *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggingMiddleware())

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

		// 菜品路由
		dishes := api.Group("/dishes")
		{
			dishes.GET("", dishHandler.List)
			dishes.GET("/:id", dishHandler.GetByID)
			dishes.POST("", middleware.AuthMiddleware(), dishHandler.Create)
			dishes.PUT("/:id", middleware.AuthMiddleware(), dishHandler.Update)
			dishes.DELETE("/:id", middleware.AuthMiddleware(), dishHandler.Delete)
			dishes.GET("/search", dishHandler.Search)
		}

		// 食材路由
		ingredients := api.Group("/ingredients")
		{
			ingredients.GET("", ingredientHandler.List)
			ingredients.POST("", middleware.AuthMiddleware(), ingredientHandler.Create)
			ingredients.PUT("/:id", middleware.AuthMiddleware(), ingredientHandler.Update)
			ingredients.DELETE("/:id", middleware.AuthMiddleware(), ingredientHandler.Delete)
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
