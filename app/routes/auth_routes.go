package routes

import (
	"classroom/app/controllers"
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // Solo para desarrollo
	}

	controller := controllers.NewAuthController(secret)
	// PATH = /api/v1/auth
	r.POST("/generate-token", controller.GenerateToken)
}
