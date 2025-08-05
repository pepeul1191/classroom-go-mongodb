package main

import (
	"classroom/app/configs"
	"classroom/app/routes"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func mountRoutes(r *gin.Engine) {
	routes.RegisterDocumentTypesRoutes(r.Group("/api/v1/document-types"))
	routes.RegisterLocationRoutes(r.Group("/api/v1/locations"))
	routes.RegisterAuthRoutes(r.Group("/api/v1/auth"))
	routes.RegisterTeacherRoutes(r.Group("/api/v1/teachers"))
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
	// cors

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000", "https://tudominio.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// routes
	mountRoutes(r)

	r.Run(":9292") // Servidor en http://localhost:8080
}
