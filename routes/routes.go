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
	user := router.Group("/api")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/profile", controllers.ProfileHandler)
		user.GET("/clients", controllers.GetClientsHandler)

	}

	// group untuk petugas (staff dan admin)
	officer := router.Group("/api/officer")
	officer.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("staff", "admin"))
	{
		officer.GET("/clients/:nosbg", controllers.GetClientDetailHandler)
		officer.POST("/ocr", controllers.OCRHandler)
		officer.POST("/submit", controllers.SubmitOCRHandler)
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
		upload := admin.Group("/uploads")
		{
			upload.GET("/submitted", controllers.GetSubmittedUploads)
			upload.GET("/:id/validate", controllers.GetUploadValidationDetail)
			upload.POST("/validate", controllers.ValidateUpload)
		}
		stats := admin.Group("/statistics")
		{
			stats.GET("/submittedUploads", controllers.GetSubmittedUploadsCount)
			stats.GET("/validatedToday", controllers.GetValidatedTodayCount)
			stats.GET("/activeOfficers", controllers.GetActiveOfficersCount)
			stats.GET("/totalSubmissions", controllers.GetTotalSubmissionsCount)
		}
	}
}
