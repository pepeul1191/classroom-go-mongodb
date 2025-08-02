package controllers

import (
	"classroom/app/models"
	"classroom/app/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocationController struct {
	Service services.LocationsService
}

func (lc *LocationController) DepartmentsFetchAll(c *gin.Context) {
	departments, err := lc.Service.FetchDepartments()
	if err != nil {
		log.Println("❌ Error al obtener departamentos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener departamentos",
			"error":   err.Error(),
		})
	}

	if len(departments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No se encontraron departametos",
			"error":   "No hay deparamentos en la colección 'locations'",
		})
	} else {
		c.JSON(http.StatusOK, departments)
	}
}

func (lc *LocationController) ProvincesFetchByDepartment(c *gin.Context) {
	deptID := c.Param("department_id")

	provinces, err := lc.Service.FetchProvincesByDepartment(deptID)
	if err != nil {
		log.Println("❌ Error al obtener departamentos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener departamentos",
			"error":   err.Error(),
		})
	}

	if len(provinces) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "El departamento no tiene provincias",
			"error":   fmt.Sprintf("No hay provincias en la colección 'locations' con el parent_id '%s'", deptID),
		})
	} else {
		c.JSON(http.StatusOK, provinces)
	}
}

func (lc *LocationController) DistrictsFetchByProvince(c *gin.Context) {
	provinceID := c.Param("province_id")

	districts, err := lc.Service.FetchDistrictsByProvince(provinceID)
	if err != nil {
		log.Println("❌ Error al obtener departamentos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener departamentos",
			"error":   err.Error(),
		})
	}

	if len(districts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "El provincia no tiene distritos",
			"error":   fmt.Sprintf("No hay distritos en la colección 'locations' con el parent_id '%s'", provinceID),
		})
	} else {
		c.JSON(http.StatusOK, districts)
	}
}

func (lc *LocationController) LocationFind(c *gin.Context) {
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

	districts, err := lc.Service.FindDistrictsByFullName(name, limitN)
	if err != nil {
		log.Println("❌ Error al buscar ubicaciones:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al buscar ubicaciones",
			"error":   err.Error(),
		})
		return
	}

	if len(districts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No hay coincidencias",
			"error":   fmt.Sprintf("No hay distritos, provincias y departamentos en la colección 'locations' que contengan la cadena '%s'", name),
		})
	} else {
		c.JSON(http.StatusOK, districts)
	}
}

func (lc *LocationController) DepartmentsCreate(c *gin.Context) {
	var input models.LocationMin

	// Intentar vincular el JSON al modelo
	if err := c.ShouldBindJSON(&input); err != nil || input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON inválido o faltan campos requeridos",
		})
		return
	}

	// Llamar al servicio para insertar
	location, err := lc.Service.InsertDepartment(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo insertar el departamento",
		})
		return
	}

	c.JSON(http.StatusCreated, location)
}

func (lc *LocationController) ProvincesCreate(c *gin.Context) {
	var input models.LocationMin
	departmentIdStr := c.Param("department_id")

	parentID, err := primitive.ObjectIDFromHex(departmentIdStr)
	if err != nil {
		fmt.Println("❌ ID inválido:", err)
		return
	}

	// Intentar vincular el JSON al modelo
	if err := c.ShouldBindJSON(&input); err != nil || input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON inválido o faltan campos requeridos",
		})
		return
	}

	// Llamar al servicio para insertar
	location, err := lc.Service.InsertProvince(input, parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo insertar el departamento",
		})
		return
	}

	c.JSON(http.StatusCreated, location)
}
