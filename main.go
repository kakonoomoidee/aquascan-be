package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"server_aquascan/config"
	"server_aquascan/routes"
)

func main() {
	// initialize database
	config.InitDatabase()

	// initialize JWT
	config.InitJWT()

	// setup router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // domain frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	routes.SetupRoutes(router)
	router.Run("0.0.0.0:8080")
}
