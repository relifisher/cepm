package services

import (
	"cepm-backend/models"
	"cepm-backend/repositories"
)

type DepartmentService struct {
	departmentRepo *repositories.DepartmentRepository
}

func NewDepartmentService(departmentRepo *repositories.DepartmentRepository) *DepartmentService {
	return &DepartmentService{departmentRepo: departmentRepo}
}

func (s *DepartmentService) CreateDepartment(department *models.Department) error {
	return s.departmentRepo.CreateDepartment(department)
}

func (s *DepartmentService) GetAllDepartments() ([]models.Department, error) {
	return s.departmentRepo.FindAllDepartments()
}
