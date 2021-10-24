package service

import (
	"red/domain"
	"red/dto"
	"red/errs"
	"red/logger"
)

//go:generate mockgen -destination=../mocks/service/mockCustomerService.go -package=service red/service CustomerService
type CustomerService interface {
	GetAllCustomer(string) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(string) (*dto.CustomerResponse, *errs.AppError)
	SaveCustomer(customer dto.CustomerRequest) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

// GetAllCustomer 找尋所有顧客的資料
func (s DefaultCustomerService) GetAllCustomer(status string) ([]dto.CustomerResponse, *errs.AppError) {
	// 轉換req裡的資料
	if status == "active" {
		status = "1"
	} else if status == "inactive" {
		status = "0"
	} else {
		status = ""
	}

	// 查詢資料
	customers, err := s.repo.FindAll(status)
	if err != nil {
		//logger.Error("failed to get all customers error")
		return nil, err
	}
	// 將資料格式轉為回傳的格式
	response := make([]dto.CustomerResponse, 0)
	for _, c := range customers {
		response = append(response, c.ToDto())
	}
	return response, nil
}

// GetCustomer 找尋特定id的顧客資料
func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}

	//將資料格式轉為回傳的格式
	response := c.ToDto()
	return &response, nil
}

// SaveCustomer 存入顧客資料
func (s DefaultCustomerService) SaveCustomer(req dto.CustomerRequest) (*dto.CustomerResponse, *errs.AppError) {
	a := domain.Customer{
		Name:        req.Name,
		City:        req.City,
		Zipcode:     req.Zipcode,
		DateOfBirth: req.DateOfBirth,
	}
	// 存入顧客資料
	c, err := s.repo.Save(a)
	if err != nil {
		logger.Error("failed to save customer")
		return nil, err
	}
	// 轉換格式為回傳需要的格式
	response := c.ToDto()
	return &response, nil
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository}
}
