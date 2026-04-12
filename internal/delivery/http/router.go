package http

import (
	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-pos/internal/delivery/http/handler"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	v1 := r.Group("/api/v1")
	{
		// Grouping khusus endpoint users
		users := v1.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.POST("/login", userHandler.Login)
			// users.GET("/:id", userHandler.GetUser)       <-- Nanti tambah di sini
			// users.PUT("/:id", userHandler.UpdateUser)    <-- Nanti tambah di sini
		}

		// Nanti jika ada fitur Product, cukup tambahkan handler-nya di parameter
		// lalu daftarkan rutenya di sini:
		// products := v1.Group("/products")
		// products.POST("/", productHandler.CreateProduct)
	}
}
