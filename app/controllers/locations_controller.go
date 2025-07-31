package controllers

import (
	"classroom/app/services"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func DepartmentsFetchAll(c *gin.Context) {
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

func ProvincesFetchByDepartment(c *gin.Context) {
	deptID := c.Param("department_id")

	results, err := services.FetchProvincesByDepartment(deptID)
	if err != nil {
		log.Println("❌ Error al obtener provincias:", err)

		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			// No se encontraron documentos (en caso de usar FindOne, por ejemplo)
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No se encontraron provincias",
				"error":   err.Error(),
			})

		default:
			// Error interno genérico
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error al obtener provincias",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, results)
}

func DistrictsFetchByProvince(c *gin.Context) {
	provinceID := c.Param("province_id")

	results, err := services.FetchDistrictsByProvince(provinceID)
	if err != nil {
		log.Println("❌ Error al obtener provincias:", err)

		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			// No se encontraron documentos (en caso de usar FindOne, por ejemplo)
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No se encontraron distritos",
				"error":   err.Error(),
			})

		default:
			// Error interno genérico
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error al obtener distritos",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, results)
}

func LocationFind(c *gin.Context) {
	name := c.Query("name")
	limit := c.Query("limit")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "El parámetro 'name' es requerido",
		})
		return
	}

	var limitN uint
	if limit == "" {
		limitN = 10
	} else {
		if _, err := fmt.Sscanf(limit, "%d", &limitN); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "limit inválido", "error": err.Error()})
			return
		}
	}

	results, err := services.FindDistrictsByFullName(name, limitN)
	if err != nil {
		log.Println("❌ Error al buscar ubicaciones:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al buscar ubicaciones",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}
