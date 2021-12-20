package service

import (
	"github.com/golang/mock/gomock"
	"red/domain"
	"red/dto"
	"red/errs"
	domain2 "red/mocks/domain"
	"reflect"
	"testing"
	"time"
)

func TestDefaultCustomerService_GetAllCustomer_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)
	customers := []domain.Customer{
		{
			Id:          1,
			UserID:      1,
			Name:        "Ivy",
			City:        "Taiwan",
			DateOfBirth: "2006-01-02",
			Status:      "1",
		},
	}

	mockRepo.EXPECT().FindAll(1).Return(customers, nil)

	want := []dto.CustomerResponse{{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "2006-01-02",
		Status:      "active",
	},
	}

	// Act
	got, _ := service.GetAllCustomers(1)
	// Assert
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Test GetAllCustomers failed:\ngot: %v\nwant: %v\n", got, want)
	}
}

func TestDefaultCustomerService_GetAllCustomer_Failed(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)

	mockRepo.EXPECT().FindAll(1).Return(nil, errs.NewUnexpectedError("Unexpected database error"))

	// Act
	_, err := service.GetAllCustomers(1)
	// Assert
	if want := errs.NewUnexpectedError("Unexpected database error"); !reflect.DeepEqual(err, want) {
		t.Errorf("Test GetAllCustomers failed:\ngot: %v\nwant: %v\n", err, want)
	}
}

func TestDefaultCustomerService_GetCustomer_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)

	customers := domain.Customer{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "2006-01-02",
		Status:      "1",
	}

	mockRepo.EXPECT().ByID("1").Return(&customers, nil)

	// Act
	got, _ := service.GetCustomer("1")

	// Assert
	want := dto.CustomerResponse{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "2006-01-02",
		Status:      "active",
	}

	if !reflect.DeepEqual(got.Name, want.Name) {
		t.Errorf("Test GetAllCustomers failed:\ngot: %v\nwant: %v\n", got, want)
	}
}

func TestDefaultCustomerService_GetCustomer_Failed(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)

	mockRepo.EXPECT().ByID("1").Return(nil, errs.NewUnexpectedError("Unexpected database error"))

	// Act
	_, err := service.GetCustomer("1")

	// Assert
	if want := errs.NewUnexpectedError("Unexpected database error"); !reflect.DeepEqual(err, want) {
		t.Errorf("Test GetAllCustomers failed:\ngot: %v\nwant: %v\n", err, want)
	}
}

func TestDefaultCustomerService_SaveCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)

	customer := domain.Customer{
		UserID:      1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "2006/01/02",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	returnCustomer := &domain.Customer{
		Id:          1,
		UserID:      1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "2006/01/02",
		Status:      "1",
	}

	mockRepo.EXPECT().Save(customer).Return(returnCustomer, nil)

	want := &dto.CustomerResponse{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "2006/01/02",
		Status:      "active",
	}

	// Act
	req := dto.CustomerRequest{
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "2006/01/02",
	}

	got, _ := service.SaveCustomer(1, req)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Test SaveCustomer: \ngot: %v\nwant %v", got, want)
	}

}

func TestDefaultCustomerService_SaveCustomer_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)

	customer := domain.Customer{
		UserID:      1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "2006/01/02",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.EXPECT().Save(customer).Return(nil, errs.NewUnexpectedError("Unexpected error from database"))

	// Act
	req := dto.CustomerRequest{
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "2006/01/02",
	}
	_, err := service.SaveCustomer(1, req)
	if want := errs.NewUnexpectedError("Unexpected error from database"); !reflect.DeepEqual(err, want) {
		t.Errorf("Test GetAllCustomers failed:\ngot: %v\nwant: %v\n", err, want)
	}

}
