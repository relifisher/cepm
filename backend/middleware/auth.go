package middleware

import (
	"net/http"
	"strings"

	"cepm-backend/models"
	"cepm-backend/services"

	"github.com/gin-gonic/gin"
)

// UserContextKey is the key to store the user in Gin context
const UserContextKey = "currentUser"

// AuthMiddleware validates JWT tokens and sets the current user in context.
func AuthMiddleware(authService services.AuthService, userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Expecting "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := authService.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token: " + err.Error()})
			c.Abort()
			return
		}

		// Fetch user from DB using UserID from claims
		user, err := userService.GetUserByID(claims.UserID) // Need to add GetUserByID to UserService
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found: " + err.Error()})
			c.Abort()
			return
		}

		c.Set(UserContextKey, user)
		c.Next()
	}
}

// GetUserFromContext retrieves the user from Gin context
func GetUserFromContext(c *gin.Context) *models.User {
	if user, exists := c.Get(UserContextKey); exists {
		if u, ok := user.(*models.User); ok {
			return u
		}
	}
	return nil
}
