package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"foodcook/internal/app/handlers"
	"foodcook/internal/app/routes"
	"foodcook/internal/infrastructure/repositories"
	"foodcook/internal/pkg/config"
	"foodcook/internal/pkg/database"

	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg := config.GetConfig()

	// 设置日志级别
	if cfg.App.Mode == "release" {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// 初始化数据库
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// 初始化数据库表和数据
	if err := database.InitTables(); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	if err := database.SeedData(); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}

	// 初始化Redis
	if err := database.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}
	defer database.CloseRedis()

	// 获取数据库连接
	db := database.GetDB()

	// 创建仓储层
	userRepo := repositories.NewMySQLUserRepository(db)
	dishRepo := repositories.NewMySQLDishRepository(db)
	ingredientRepo := repositories.NewMySQLIngredientRepository(db)
	mealRecordRepo := repositories.NewMySQLMealRecordRepository(db)
	categoryRepo := repositories.NewMySQLCategoryRepository(db)

	// 创建处理器
	authHandler := handlers.NewAuthHandler(userRepo)
	dishHandler := handlers.NewDishHandler(dishRepo)
	ingredientHandler := handlers.NewIngredientHandler(ingredientRepo)
	mealRecordHandler := handlers.NewMealRecordHandler(mealRecordRepo, dishRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)

	// 设置路由
	r := routes.SetupRoutes(authHandler, dishHandler, ingredientHandler, mealRecordHandler, categoryHandler)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", cfg.App.Port),
		Handler: r,
	}

	// 启动服务器
	go func() {
		logrus.Infof("Server starting on port %d", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}

	logrus.Info("Server exited")
}
