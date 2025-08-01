package mocks

import (
	"classroom/app/models"
)

type MockDepartmentService struct {
	Departments []models.LocationMin
	Err         error
}

func (m *MockDepartmentService) FetchDepartments() ([]models.LocationMin, error) {
	return m.Departments, m.Err
}
