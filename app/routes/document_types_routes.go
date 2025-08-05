package routes

import (
	"classroom/app/controllers"
	"classroom/app/services"

	"github.com/gin-gonic/gin"
)

func RegisterDocumentTypesRoutes(r *gin.RouterGroup) {
	service := services.NewDocumentTypesService()
	controller := &controllers.DocumentTypeController{Service: service}
	// PATH = /api/v1/document_types
	r.GET("", controller.DocumentTypeFetchAll)
	r.GET("/:_id", controller.DocumentTypeFetchOne)
}
