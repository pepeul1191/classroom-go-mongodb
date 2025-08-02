package test

import (
	"bytes"
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

func TestDepartmentsCreate_Success(t *testing.T) {
	mock := &mocks.MockLocationsService{}

	router := gin.Default()
	controller := &controllers.LocationController{Service: mock}
	router.POST("/api/v1/locations/departments", controller.DepartmentsCreate)

	body := `{"name":"Lima"}`
	req, _ := http.NewRequest("POST", "/api/v1/locations/departments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var loc models.Location
	json.Unmarshal(w.Body.Bytes(), &loc)
	println("1 ++++++++++++++++++++++++++++++++++++++++")
	println(loc.ID.Hex())
	println("2 ++++++++++++++++++++++++++++++++++++++++")
	assert.Equal(t, "Lima", loc.Name)
	assert.Equal(t, "department", loc.Type)
	// Validar que el ID no esté vacío (puede ser cualquier ObjectID válido)
	assert.False(t, loc.ID.IsZero(), "ID debe estar presente")

	// Validar que las fechas estén inicializadas
	assert.False(t, loc.Created.IsZero(), "Created no debe estar vacío")
	assert.False(t, loc.Updated.IsZero(), "Updated no debe estar vacío")

	// También puedes imprimir si quieres ver los valores:
	t.Logf("ID: %s\nCreated: %s\nUpdated: %s", loc.ID.Hex(), loc.Created.String(), loc.Updated.String())
}

func TestDepartmentsCreate_MissingName(t *testing.T) {
	mock := &mocks.MockLocationsService{}

	router := gin.Default()
	controller := &controllers.LocationController{Service: mock}
	router.POST("/api/v1/locations/departments", controller.DepartmentsCreate)

	// Body con campo "name" vacío
	body := `{"name":""}`
	req, _ := http.NewRequest("POST", "/api/v1/locations/departments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "JSON inválido")
}
