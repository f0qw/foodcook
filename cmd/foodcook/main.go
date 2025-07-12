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

	// 初始化Redis
	if err := database.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}
	defer database.CloseRedis()

	// 创建处理器（这里需要实现具体的仓储层）
	// TODO: 实现具体的仓储层
	authHandler := &handlers.AuthHandler{}
	dishHandler := &handlers.DishHandler{}
	ingredientHandler := &handlers.IngredientHandler{}
	mealRecordHandler := &handlers.MealRecordHandler{}

	// 设置路由
	r := routes.SetupRoutes(authHandler, dishHandler, ingredientHandler, mealRecordHandler)

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
