package mocks

import (
	"classroom/app/models"
)

type MockLocationsService struct {
	Departments      []models.LocationMin
	Err              error
	ExpectedLocation *models.Location
}

// respuesta esperada
func (m *MockLocationsService) FetchDepartments() ([]models.LocationMin, error) {
	return m.Departments, m.Err
}

func (m *MockLocationsService) InsertDepartment(dep models.LocationMin) (*models.Location, error) {
	return m.ExpectedLocation, m.Err
}

// Los otros métodos pueden dejarse vacíos si no los necesitas para este test
func (m *MockLocationsService) FetchProvincesByDepartment(string) ([]models.LocationMin, error) {
	return nil, nil
}

func (m *MockLocationsService) FetchDistrictsByProvince(string) ([]models.LocationMin, error) {
	return nil, nil
}

func (m *MockLocationsService) FindDistrictsByFullName(string, uint) ([]models.LocationResult, error) {
	return nil, nil
}
