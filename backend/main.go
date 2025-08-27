package main

import (
	"fmt"
	"log"

	"cepm-backend/config"
	"cepm-backend/database"
	"cepm-backend/models"
	"cepm-backend/repositories"
	"cepm-backend/router"
	"cepm-backend/services"
	"cepm-backend/wechat"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	database.Init(&cfg.Database)

	// Auto-migrate the schema
	// This will create tables, columns, and foreign keys. 
	// It's safe to run every time, as it will only add missing things.
	models.AutoMigrate(database.DB)

	// Seed database with mock data (for development)
	database.SeedData(database.DB)

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)

	departmentRepo := repositories.NewDepartmentRepository(database.DB)
	departmentService := services.NewDepartmentService(departmentRepo)

	systemSettingRepo := repositories.NewSystemSettingRepository(database.DB)
	systemSettingService := services.NewSystemSettingService(systemSettingRepo)

	// Initialize WeChat Client
	wechatClient := wechat.NewWechatClient(&cfg.Wechat)

	// Initialize Auth Service
	authService := services.NewAuthService(userRepo, wechatClient, &cfg.JWT)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Setup router
	r := router.SetupRouter(userService, departmentService, systemSettingService, authService)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
