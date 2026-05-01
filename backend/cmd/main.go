package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"recruit/database"
	"recruit/models"
	"recruit/routes"
	"recruit/utils"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// 连接数据库
	database.Connect()

	// 自动迁移数据库
	db := database.GetDB()
	err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Bid{},
		&models.Review{},
		&models.Favorite{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 创建默认管理员账户
	createDefaultAdmin(db)

	// 设置路由
	r := routes.SetupRouter()

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func createDefaultAdmin(db *gorm.DB) {
	var admin models.User
	if err := db.Where("role = ?", models.RoleAdmin).First(&admin).Error; err == nil {
		log.Println("Admin user already exists")
		return
	} else if err != gorm.ErrRecordNotFound {
		log.Println("Warning: Failed to check admin user:", err)
		return
	}

	// 创建默认管理员
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Println("Warning: Failed to hash admin password:", err)
		return
	}

	admin = models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: hashedPassword,
		Role:     models.RoleAdmin,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Println("Warning: Failed to create admin user:", err)
		return
	}

	log.Println("Default admin user created:")
	log.Println("  Email: admin@example.com")
	log.Println("  Password: admin123")
}
