package test

import (
	"classroom/app/controllers"
	"classroom/app/models"
	"classroom/tests/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// instancia de prueba del ruteador
func setupRouterWithMock(service *mocks.MockLocationsService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	// controlador que quiero probar
	controller := &controllers.LocationController{
		Service: service,
	}
	// ruta a probar
	r.GET("/api/v1/locations/departments", controller.DepartmentsFetchAll)
	return r
}

func TestDepartmentsFetchAll_Success(t *testing.T) {
	mockService := &mocks.MockLocationsService{
		Departments: []models.LocationMin{
			{ID: primitive.NewObjectID(), Name: "Lima"},
			{ID: primitive.NewObjectID(), Name: "Cusco"},
		},
	}

	router := setupRouterWithMock(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/locations/departments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body []models.LocationMin
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Len(t, body, 2)
}

func TestDepartmentsFetchAll_NotFound(t *testing.T) {
	mockService := &mocks.MockLocationsService{
		Departments: []models.LocationMin{}, // vacío
	}

	router := setupRouterWithMock(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/locations/departments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDepartmentsFetchAll_Error(t *testing.T) {
	mockService := &mocks.MockLocationsService{
		Err: errors.New("DB error"),
	}

	router := setupRouterWithMock(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/locations/departments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
