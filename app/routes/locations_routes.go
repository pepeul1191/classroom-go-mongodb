package routes

import (
	"classroom/app/controllers"
	"classroom/app/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterLocationRoutes(r *gin.RouterGroup) {
	service := services.NewLocationsService()
	controller := &controllers.LocationController{Service: service}
	// PATH = /api/v1/locations
	fmt.Println("-------------------------------------------------------")
	r.GET("/departments", controller.DepartmentsFetchAll)
	r.POST("/departments", controller.SaveDepartments)
	r.GET("/departments/:department_id/provinces", controller.ProvincesFetchByDepartment)
	r.POST("/departments/:department_id/provinces", controller.ProvincesCreate)
	r.GET("/provinces/:province_id/districts", controller.DistrictsFetchByProvince)
	r.GET("/find", controller.LocationFind)
}
