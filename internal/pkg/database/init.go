package database

import (
	"foodcook/internal/domain/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// InitTables 初始化数据库表
func InitTables() error {
	if DB == nil {
		return nil
	}

	// 自动迁移表结构
	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Dish{},
		&models.Ingredient{},
		&models.DishIngredient{},
		&models.MealRecord{},
		&models.MealRecordDish{},
	)
	if err != nil {
		return err
	}

	log.Println("Database tables initialized successfully")
	return nil
}

// SeedData 插入初始数据
func SeedData() error {
	if DB == nil {
		return nil
	}

	// 检查是否已有数据
	var count int64
	DB.Model(&models.Category{}).Count(&count)
	if count > 0 {
		log.Println("Database already has data, skipping seed")
		return nil
	}

	// 创建 root 用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("root123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return err
	}

	rootUser := models.User{
		Username:     "root",
		Email:        "root@foodcook.com",
		PasswordHash: string(hashedPassword),
		Role:         models.RoleRoot,
	}

	if err := DB.Create(&rootUser).Error; err != nil {
		log.Printf("Failed to create root user: %v", err)
	}

	// 插入默认分类
	categories := []models.Category{
		{Name: "川菜"},
		{Name: "粤菜"},
		{Name: "湘菜"},
		{Name: "鲁菜"},
		{Name: "苏菜"},
		{Name: "浙菜"},
		{Name: "闽菜"},
		{Name: "徽菜"},
		{Name: "东北菜"},
		{Name: "西北菜"},
		{Name: "西南菜"},
		{Name: "其他"},
	}

	for _, category := range categories {
		if err := DB.Create(&category).Error; err != nil {
			log.Printf("Failed to create category %s: %v", category.Name, err)
		}
	}

	// 插入示例菜品
	dishes := []models.Dish{
		{
			Name:        "麻婆豆腐",
			Description: "四川传统名菜，麻辣鲜香",
			Price:       28.00,
			CategoryID:  &[]uint{1}[0], // 川菜
			ImageURL:    "https://example.com/mapo-tofu.jpg",
		},
		{
			Name:        "白切鸡",
			Description: "广东名菜，皮爽肉嫩",
			Price:       45.00,
			CategoryID:  &[]uint{2}[0], // 粤菜
			ImageURL:    "https://example.com/white-cut-chicken.jpg",
		},
		{
			Name:        "剁椒鱼头",
			Description: "湖南特色菜，酸辣开胃",
			Price:       68.00,
			CategoryID:  &[]uint{3}[0], // 湘菜
			ImageURL:    "https://example.com/chopped-pepper-fish-head.jpg",
		},
		{
			Name:        "糖醋里脊",
			Description: "经典家常菜，酸甜可口",
			Price:       32.00,
			CategoryID:  &[]uint{4}[0], // 鲁菜
			ImageURL:    "https://example.com/sweet-sour-pork.jpg",
		},
	}

	for _, dish := range dishes {
		if err := DB.Create(&dish).Error; err != nil {
			log.Printf("Failed to create dish %s: %v", dish.Name, err)
		}
	}

	// 插入示例食材
	ingredients := []models.Ingredient{
		{Name: "豆腐", Unit: "块", Price: 3.00},
		{Name: "猪肉", Unit: "斤", Price: 25.00},
		{Name: "鸡肉", Unit: "斤", Price: 18.00},
		{Name: "鱼头", Unit: "个", Price: 15.00},
		{Name: "辣椒", Unit: "斤", Price: 8.00},
		{Name: "葱", Unit: "斤", Price: 5.00},
		{Name: "姜", Unit: "斤", Price: 12.00},
		{Name: "蒜", Unit: "斤", Price: 6.00},
	}

	for _, ingredient := range ingredients {
		if err := DB.Create(&ingredient).Error; err != nil {
			log.Printf("Failed to create ingredient %s: %v", ingredient.Name, err)
		}
	}

	log.Println("Database seeded successfully")
	return nil
}
