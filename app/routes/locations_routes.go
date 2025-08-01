package routes

import (
	"classroom/app/controllers"
	"classroom/app/services"

	"github.com/gin-gonic/gin"
)

func RegisterLocationRoutes(r *gin.RouterGroup) {
	service := services.NewLocationsService()
	controller := &controllers.LocationController{Service: service}

	r.GET("/departments", controller.DepartmentsFetchAll)
	r.GET("/departments/:department_id/provinces", controller.ProvincesFetchByDepartment)
	r.GET("/provinces/:province_id/districts", controller.DistrictsFetchByProvince)
	r.GET("/find", controller.LocationFind)
}
