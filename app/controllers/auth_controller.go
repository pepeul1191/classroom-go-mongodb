// controllers/auth_controller.go
package controllers

import (
	"classroom/app/models"
	"classroom/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthController struct {
	JWTSecret string
}

func NewAuthController(secret string) *AuthController {
	return &AuthController{JWTSecret: secret}
}

func (ac *AuthController) GenerateToken(c *gin.Context) {
	var req models.TokenRequest

	// Validar la estructura básica con Gin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":  "Datos de entrada inválidos",
			"detail": err.Error(),
		})
		return
	}

	// Validación adicional personalizada si es necesaria
	if err := utils.ValidateRoles(req.Roles); err != nil {
		c.JSON(400, gin.H{
			"error":  "Validación de roles falló",
			"detail": err.Error(),
		})
		return
	}

	// Simulación: obtener persona (normalmente desde BD)
	// Convertir el string ID a ObjectID
	objID, err := primitive.ObjectIDFromHex("688bc5a09cb60ad40cbe61dc")
	if err != nil {
		c.JSON(400, gin.H{
			"error":  "Al capturar el usuario",
			"detail": err.Error(),
		})
		return
	}
	person := models.Person{
		ID:        objID,
		Names:     "Pepe",
		LastNames: "Valdivia",
		ImageURL:  "user-default.png",
	}

	// Crear token JWT
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   "your-app.com",
		"aud":   "your-client-id",
		"sub":   "user@example.com",
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
		"user":  req.User.Username,
		"roles": req.Roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(ac.JWTSecret))
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al generar token"})
		return
	}

	// Respuesta
	response := models.TokenResponse{
		Person: person,
		Token:  signedToken,
	}

	c.JSON(200, response)
}
