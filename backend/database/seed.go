package database

import (
	"log"

	"cepm-backend/models"
	"gorm.io/gorm"
)

// SeedData populates the database with initial mock data for development.
func SeedData(db *gorm.DB) {
	// Check if data has already been seeded by checking for a specific user.
	var userCount int64
	db.Model(&models.User{}).Where("email = ?", "manager@example.com").Count(&userCount)
	if userCount > 0 {
		log.Println("Database already contains seed data. Skipping.")
		return
	}

	log.Println("Seeding database with mock data...")

	// 1. Create Roles with error checking
	roles := []models.Role{{Name: "员工"}, {Name: "组长"}, {Name: "中心负责人"}, {Name: "人事"}, {Name: "管理员"}}
	for _, r := range roles {
		if err := db.Where(models.Role{Name: r.Name}).FirstOrCreate(&r).Error; err != nil {
			log.Fatalf("failed to seed role %s: %v", r.Name, err)
		}
	}

	// 2. Create a Manager
	var teamLeadRole models.Role
	db.Where("name = ?", "组长").First(&teamLeadRole)
	manager := models.User{Name: "王经理", Email: "manager@example.com", WechatUserid: "wechat_manager", RoleID: &teamLeadRole.ID}
	if err := db.Where(models.User{Email: manager.Email}).FirstOrCreate(&manager).Error; err != nil {
		log.Fatalf("failed to seed manager: %v", err)
	}

	// 3. Create Employees
	var employeeRole models.Role
	db.Where("name = ?", "员工").First(&employeeRole)
	employeeLisi := models.User{Name: "李四", Email: "lisi@example.com", WechatUserid: "wechat_lisi", RoleID: &employeeRole.ID, ManagerID: &manager.ID}
	if err := db.Where(models.User{Email: employeeLisi.Email}).FirstOrCreate(&employeeLisi).Error; err != nil {
		log.Fatalf("failed to seed employee lisi: %v", err)
	}

	employeeZhaowu := models.User{Name: "赵五", Email: "zhaowu@example.com", WechatUserid: "wechat_zhaowu", RoleID: &employeeRole.ID, ManagerID: &manager.ID}
	if err := db.Where(models.User{Email: employeeZhaowu.Email}).FirstOrCreate(&employeeZhaowu).Error; err != nil {
		log.Fatalf("failed to seed employee zhaowu: %v", err)
	}

	// 4. Create a sample performance review for Li Si
	log.Println("Seeding performance review for Li Si...")
	review := models.PerformanceReview{
		UserID: employeeLisi.ID,
		Period: "2025-07",
		Status: "待打分",
		Items: []models.PerformanceItem{
			{Category: "工作业绩", Title: "完成V2.0模块开发", Weight: 50, Target: "V2.0版本按时上线"},
			{Category: "工作业绩", Title: "修复线上BUG", Weight: 30, Target: "BUG数量减少50%"},
			{Category: "大模型", Title: "大模型使用能力", Weight: 10, Target: "在日常工作中有效使用AI工具"},
			{Category: "价值观", Title: "价值观践行", Weight: 10, Target: "积极参与团队分享"},
		},
	}
	if err := db.Create(&review).Error; err != nil {
		log.Fatalf("failed to seed performance review: %v", err)
	}

	// 5. Create an Admin User
	var adminRole models.Role
	db.Where("name = ?", "管理员").First(&adminRole)
	adminUser := models.User{Name: "管理员", Email: "admin@example.com", WechatUserid: "wechat_admin", RoleID: &adminRole.ID}
	if err := db.Where(models.User{Email: adminUser.Email}).FirstOrCreate(&adminUser).Error; err != nil {
		log.Fatalf("failed to seed admin user: %v", err)
	}

	log.Println("Database seeding completed successfully.")
}
