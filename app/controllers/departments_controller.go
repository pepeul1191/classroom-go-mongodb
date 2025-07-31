package controllers

import (
	"classroom/app/services"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func DepartmentFetchAll(c *gin.Context) {
	results, err := services.FetchDepartments()
	if err != nil {
		log.Println("❌ Error al obtener departamentos:", err)

		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			// No se encontraron documentos (en caso de usar FindOne, por ejemplo)
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No se encontraron departamentos",
				"error":   err.Error(),
			})

		default:
			// Error interno genérico
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error al obtener departamentos",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, results)
}
