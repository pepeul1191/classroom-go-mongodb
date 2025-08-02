package main

import (
	"classroom/app/configs"
	"classroom/app/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func mountRoutes(r *gin.Engine) {
	routes.RegisterLocationRoutes(r.Group("/api/v1/locations"))
}

func main() {
	// mongodb
	err := configs.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	// app
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("X-Powered-By", "Gin")
		c.Header("Server", "Ubuntu")
		c.Next()
	})
	r.Static("/static", "./public")

	mountRoutes(r)

	r.Run(":8080") // Servidor en http://localhost:8080
}
