package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"classroom/app/controllers"
	"classroom/app/models"
	"classroom/test/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouterWithController(ctrl *controllers.LocationController) *gin.Engine {
	r := gin.Default()
	r.GET("/departments", ctrl.DepartmentsFetchAll)
	return r
}

func TestDepartmentsFetchAll_Success(t *testing.T) {
	mock := &mocks.MockDepartmentService{
		Departments: []models.LocationMin{
			{ID: "1", Name: "Lima"},
		},
	}

	ctrl := &controllers.LocationController{Service: mock}
	router := setupRouterWithController(ctrl)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.LocationMin
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Lima", response[0].Name)
}
