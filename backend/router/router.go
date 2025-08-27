package router

import (
	"net/http"
	"time"

	"cepm-backend/api"
	"cepm-backend/middleware"
	"cepm-backend/repositories"
	"cepm-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func SetupRouter(userService *services.UserService, departmentService *services.DepartmentService, systemSettingService *services.SystemSettingService, authService services.AuthService) *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3100"}, // Allow your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-Email"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Dependency Injection
	performanceReviewRepo := repositories.NewPerformanceReviewRepository()
	performanceReviewService := services.NewPerformanceReviewService(performanceReviewRepo)
	performanceReviewHandler := api.NewPerformanceReviewHandler(performanceReviewService)
	adminHandler := api.NewAdminHandler(userService, departmentService, systemSettingService)
	authHandler := api.NewAuthHandler(authService)

	// API v1 group
	apiV1 := r.Group("/api/v1")
	{
		// Public auth routes
		apiV1.GET("/wechat/login", authHandler.WechatLogin)

		// Apply AuthMiddleware to all subsequent routes in this group
		apiV1.Use(middleware.AuthMiddleware(authService, userService))

		// Performance Review routes
		reviews := apiV1.Group("/reviews")
		{
			// Use "" for the group's root path, e.g., /api/v1/reviews
			reviews.POST("", performanceReviewHandler.CreatePerformanceReview)
			reviews.GET("", performanceReviewHandler.ListUserReviews)
			// New route for getting a review by user and period
			reviews.GET("/by-period", performanceReviewHandler.GetPerformanceReviewByPeriod)
			reviews.GET("/all-submitted", performanceReviewHandler.ListAllSubmittedReviews) // New route for HR role
			reviews.GET("/all-by-period", performanceReviewHandler.ListAllReviewsByPeriod) // New route for HR to view all reviews by period

			// Routes with path parameters
			reviews.GET("/:id", performanceReviewHandler.GetPerformanceReview)
			reviews.PUT("/:id", performanceReviewHandler.UpdatePerformanceReview)
			reviews.POST("/:id/score", performanceReviewHandler.ScorePerformanceReview)
			reviews.POST("/:id/submit", performanceReviewHandler.SubmitPerformanceReview)
			reviews.POST("/:id/approve", performanceReviewHandler.ApprovePerformanceReview)
			reviews.POST("/:id/reject", performanceReviewHandler.RejectPerformanceReview)
		}

		// Team-related routes
		team := apiV1.Group("/team")
		{
			team.GET("/reviews", performanceReviewHandler.ListTeamReviews)
		}

		// Admin routes
		admin := apiV1.Group("/admin")
		admin.Use(middleware.RequireRole("管理员"))
		{
			admin.GET("/users", adminHandler.GetUsers)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.POST("/departments", adminHandler.CreateDepartment)
			admin.GET("/departments", adminHandler.GetDepartments)
			admin.GET("/roles", adminHandler.GetRoles)
			admin.PUT("/settings", adminHandler.UpdateSystemSetting)
		}
	}

	return r
}