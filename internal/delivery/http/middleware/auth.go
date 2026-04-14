package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-pos/pkg/response"
	"github.com/valenrio66/be-pos/pkg/utils"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Access denied", "Authorization header not found")
			c.Abort()
			return
		}

		claims, err := utils.VerifyToken(authHeader, secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Access denied", "The token is invalid or has expired")
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Access denied", "The token value is invalid")
			c.Abort()
			return
		}

		c.Set("user_id", int64(userIDFloat))
		c.Set("user_role", claims["role"])

		c.Next()
	}
}
