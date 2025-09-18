package routes

import (
	"github.com/gin-gonic/gin"

	"server_aquascan/controllers"
	"server_aquascan/middleware"
)

func SetupRoutes(router *gin.Engine) {
	// Endpoint yang tidak memerlukan autentikasi
	public := router.Group("/api")
	{
		public.POST("/login", controllers.LoginHandler)
	}
	// Endpoint yang memerlukan autentikasi
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", controllers.ProfileHandler)
		protected.POST("/upload", controllers.UploadHandler)
	}
}
