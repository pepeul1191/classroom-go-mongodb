package routes

import (
	"classroom/app/controllers"
	"classroom/app/services"

	"github.com/gin-gonic/gin"
)

func RegisterTeacherRoutes(r *gin.RouterGroup) {
	service := services.NewTeachersService()
	controller := &controllers.TeacherController{Service: service}
	// PATH = /api/v1/teachers
	r.POST("", controller.SaveTeacher)
	r.PUT("/:_id", controller.UpdateTeacher)
}
