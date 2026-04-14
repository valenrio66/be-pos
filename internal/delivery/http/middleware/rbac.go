package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-pos/pkg/response"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleVal, exists := c.Get("user_role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Unauthorized access", "User role not found in context")
			c.Abort()
			return
		}

		userRole := userRoleVal.(string)
		isAllowed := false

		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			response.Error(c, http.StatusForbidden, "Access denied", "You do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}
