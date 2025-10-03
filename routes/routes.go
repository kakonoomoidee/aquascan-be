package routes

import (
	"github.com/gin-gonic/gin"

	"server_aquascan/controllers"
	"server_aquascan/middleware"
)

func SetupRoutes(router *gin.Engine) {
	// Endpoint public
	public := router.Group("/api")
	{
		public.POST("/login", controllers.LoginHandler)
	}

	// group untuk user login biasa
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", controllers.ProfileHandler)
		protected.POST("/upload", controllers.UploadHandler)

		protected.GET("/clients", controllers.GetClientsHandler)
		protected.GET("/clients/:nosbg", controllers.GetClientDetailHandler)
	}

	// group khusus admin
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		users := admin.Group("/users")
		{
			users.GET("/all", controllers.GetAllUsersHandler)
			users.PUT("/:id", controllers.UpdateUserHandler)
			users.DELETE("/:id", controllers.DeleteUserHandler)
			users.POST("/add", controllers.RegisterHandler)
		}
		client := admin.Group("/clients")
		{
			client.GET("/:nosbg", controllers.GetMoreClientDetailHandler)
		}
	}
}
