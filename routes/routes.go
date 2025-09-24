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

		protected.GET("/clients", controllers.GetClients)
		protected.GET("/clients/:nosbg", controllers.GetClientDetail)
	}

	// group khusus admin
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		admin.POST("/regisuser", controllers.RegisterHandler)

		users := admin.Group("/users")
		{
			users.GET("/", controllers.GetAllUsersHandler)
			users.PUT("/:id", controllers.UpdateUserHandler)
			users.DELETE("/:id", controllers.DeleteUserHandler)
		}
	}
}
