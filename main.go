package main

import (
	"classroom/app/configs"
	"classroom/app/controllers"
	"classroom/app/services"
	"log"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {
	service := services.NewLocationsService()
	locationController := &controllers.LocationController{
		Service: service,
	}

	r.GET("/api/v1/locations/departments", locationController.DepartmentsFetchAll)
	//r.POST("/api/v1/locations/departments", locationController.DepartmentsSave)
	//r.PUT("/api/v1/locations/departments/:department_id", locationController.DepartmentsUpdate)
	//r.DELETE("/api/v1/locations/departments/:department_id", locationController.DepartmentsDelete)
	r.GET("/api/v1/locations/departments/:department_id/provinces", locationController.ProvincesFetchByDepartment)
	r.GET("/api/v1/locations/provinces/:province_id/districts", locationController.DistrictsFetchByProvince)
	r.GET("/api/v1/locations/find", locationController.LocationFind)
}

func main() {
	// mongodb
	err := configs.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	// app
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("X-Powered-By", "Gin")
		c.Header("Server", "Ubuntu")
		c.Next()
	})
	r.Static("/static", "./public")

	setupRoutes(r)

	r.Run(":8080") // Servidor en http://localhost:8080
}
