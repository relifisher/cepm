package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole is a middleware that checks if the authenticated user has one of the required roles.
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUserFromContext(c)
		if user == nil || user.Role.Name == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User or role information not available"})
			c.Abort()
			return
		}

		userRole := user.Role.Name
		
		// Check if the user's role is in the list of required roles
		for _, role := range requiredRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient role permissions"})
		c.Abort()
	}
}
