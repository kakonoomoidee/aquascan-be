package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"server_aquascan/routes"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	routes.SetupRoutes(router)
	router.Run("0.0.0.0:8080")
}
