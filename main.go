package main

import (
	"classroom/app/configs"
	"classroom/app/controllers"
	"log"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {
	r.GET("/api/v1/locations/departments", controllers.DepartmentsFetchAll)
	r.POST("/api/v1/locations/departments", controllers.DepartmentsSave)
	r.PUT("/api/v1/locations/departments/:department_id", controllers.DepartmentsUpdate)
	r.DELETE("/api/v1/locations/departments/:department_id", controllers.DepartmentsDelete)
	r.GET("/api/v1/locations/departments/:department_id/provinces", controllers.ProvincesFetchByDepartment)
	r.GET("/api/v1/locations/provinces/:province_id/districts", controllers.DistrictsFetchByProvince)
	r.GET("/api/v1/locations/find", controllers.LocationFind)
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
