package controllers

import (
	"classroom/app/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentTypeController struct {
	Service services.DocumentTypesService
}

func (lc *DocumentTypeController) DocumentTypeFetchAll(c *gin.Context) {
	gin.DefaultWriter.Write([]byte("Mensaje de depuración\n"))

	documentTypes, err := lc.Service.FetchAll()
	if err != nil {
		log.Println("❌ Error al obtener departamentos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener departamentos",
			"error":   err.Error(),
		})
	}

	if len(documentTypes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No se encontraron tipos de documentos",
			"error":   "No hay tipos de documentos en la colección 'document_types'",
		})
	} else {
		c.JSON(http.StatusOK, documentTypes)
	}
}

func (lc *DocumentTypeController) DocumentTypeFetchOne(c *gin.Context) {
	strID := c.Param("_id")
	documentTypeId, err := primitive.ObjectIDFromHex(strID)
	var ID *primitive.ObjectID = &documentTypeId
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "DepartmentID inválido", "deparmentId": strID})
		return
	}

	documentType, err := lc.Service.FetchOne(ID)
	if err != nil {
		log.Println("❌ Error al obtener tipo de documento:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error al obtener tipo de documento",
			"error":   err.Error(),
		})
		return
	}

	if documentType == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "El tipo de documento no encontrado",
			"error":   fmt.Sprintf("No hay un tipo de documento en la colección 'document_types' con el _id '%s'", ID),
		})
		return
	} else {
		c.JSON(http.StatusOK, documentType)
	}
}
