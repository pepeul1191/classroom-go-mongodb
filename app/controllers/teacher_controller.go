package controllers

import (
	"classroom/app/models"
	"classroom/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeacherController struct {
	Service services.TeachersService
}

func (tc *TeacherController) SaveTeacher(c *gin.Context) {
	var input models.TeacherCreateRequest

	// 1. Parsear y validar el input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	teacher, err := tc.Service.CreateTeacher(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, teacher)
}

func (tc *TeacherController) UpdateTeacher(c *gin.Context) {
	IdStr := c.Param("_id")

	ID, err := primitive.ObjectIDFromHex(IdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "_id de docente a editar no se pudo parsear",
			"details": err.Error(),
		})
		return
	}

	var input models.TeacherCreateRequest

	// 1. Parsear y validar el input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	teacher, err := tc.Service.UpdateTeacher(c.Request.Context(), &input, ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, teacher)
}
