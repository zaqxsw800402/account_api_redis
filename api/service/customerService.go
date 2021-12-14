package service

import (
	"net/http"
	"red/domain"
	"red/dto"
	"red/errs"
	"time"
)

//go:generate mockgen -destination=../mocks/service/mockCustomerService.go -package=service red/service CustomerService
type CustomerService interface {
	GetAllCustomer(int) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(string) (*dto.CustomerResponse, *errs.AppError)
	SaveCustomer(userID int, req dto.CustomerRequest) (*dto.CustomerResponse, *errs.AppError)
	NewCustomer(dto.CustomerRequest) (*dto.CustomerResponse, *errs.AppError)
	DeleteCustomer(string) *errs.AppError
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository}
}

// GetAllCustomer 找尋所有顧客的資料
func (s DefaultCustomerService) GetAllCustomer(userID int) ([]dto.CustomerResponse, *errs.AppError) {
	// 轉換req裡的資料
	//status = transformStatus(status)

	// 查詢資料
	customers, err := s.repo.FindAll(userID)
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

func transformStatus(status string) string {
	if status == "active" {
		status = "1"
	} else if status == "inactive" {
		status = "0"
	} else {
		status = ""
	}
	return status
}

// GetCustomer 找尋特定id的顧客資料
func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}

	//將資料格式轉為回傳的格式
	response := c.ToDto()
	return &response, nil
}

// SaveCustomer 存入顧客資料
func (s DefaultCustomerService) SaveCustomer(userID int, req dto.CustomerRequest) (*dto.CustomerResponse, *errs.AppError) {
	customer := domain.Customer{
		UserID: uint(userID),
		Name:   req.Name,
		City:   req.City,
		//Zipcode:     req.Zipcode,
		DateOfBirth: req.DateOfBirth,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if !customer.IsValid() {
		return nil, &errs.AppError{Code: http.StatusBadRequest, Message: "check your birthdate"}
	}

	// 存入顧客資料
	c, err := s.repo.Save(customer)
	if err != nil {
		return nil, err
	}
	// 轉換格式為回傳需要的格式
	response := c.ToDto()
	return &response, nil
}

func (s DefaultCustomerService) NewCustomer(req dto.CustomerRequest) (*dto.CustomerResponse, *errs.AppError) {
	customer := domain.Customer{
		//Id:          uint(req.Id),
		Name: req.Name,
		City: req.City,
		//Zipcode:     req.Zipcode,
		DateOfBirth: req.DateOfBirth,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	c, err := s.repo.Update(customer)
	if err != nil {
		return nil, err
	}

	response := c.ToDto()
	return &response, nil
}

func (s DefaultCustomerService) DeleteCustomer(id string) *errs.AppError {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
