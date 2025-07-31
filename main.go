package main

import (
	"classroom/app/configs"
	"classroom/app/controllers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	err := configs.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// Servir archivos estáticos (JS, CSS, HTML, imágenes, etc.)
	r.Static("/static", "./public")

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"mensaje": "¡Hola Mundo!",
		})
	})

	r.GET("/api/v1/locations/departments", controllers.DepartmentsFetchAll)
	r.GET("/api/v1/locations/departments/:department_id/provinces", controllers.ProvincesFetchByDepartment)
	r.GET("/api/v1/locations/provinces/:province_id/districts", controllers.DistrictsFetchByProvince)
	r.GET("/api/v1/locations/find", controllers.LocationFind)

	r.Run(":8080") // Servidor en http://localhost:8080
}
