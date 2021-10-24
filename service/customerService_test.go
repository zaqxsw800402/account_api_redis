package service

import (
	"github.com/golang/mock/gomock"
	"red/domain"
	"red/dto"
	"red/errs"
	domain2 "red/mocks/domain"
	"reflect"
	"testing"
)

func TestDefaultCustomerService_GetAllCustomer_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)
	customers := []domain.Customer{
		{
			Id:          1,
			Name:        "Ivy",
			City:        "Taiwan",
			Zipcode:     "23",
			DateOfBirth: "2006-01-02",
			Status:      "1",
		},
	}

	mockRepo.EXPECT().FindAll("1").Return(customers, nil)

	want := []dto.CustomerResponse{{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "23",
		DateOfBirth: "2006-01-02",
		Status:      "active",
		Accounts:    []domain.Account{},
	},
	}

	// Act
	got, _ := service.GetAllCustomer("active")
	// Assert
	if !reflect.DeepEqual(got[0].Name, want[0].Name) {
		t.Errorf("Test GetAllCustomer failed:\ngot: %v\nwant: %v\n", got, want)
	}
}

func TestDefaultCustomerService_GetAllCustomer_Failed(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)

	mockRepo.EXPECT().FindAll("1").Return(nil, errs.NewUnexpectedError("Unexpected database error"))

	// Act
	_, err := service.GetAllCustomer("active")
	// Assert
	if want := errs.NewUnexpectedError("Unexpected database error"); !reflect.DeepEqual(err, want) {
		t.Errorf("Test GetAllCustomer failed:\ngot: %v\nwant: %v\n", err, want)
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
		Zipcode:     "23",
		DateOfBirth: "2006-01-02",
		Status:      "1",
	}

	mockRepo.EXPECT().ById("1").Return(&customers, nil)

	// Act
	got, _ := service.GetCustomer("1")

	// Assert
	want := dto.CustomerResponse{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "23",
		DateOfBirth: "2006-01-02",
		Status:      "active",
		Accounts:    []domain.Account{},
	}

	if !reflect.DeepEqual(got.Name, want.Name) {
		t.Errorf("Test GetAllCustomer failed:\ngot: %v\nwant: %v\n", got, want)
	}
}

func TestDefaultCustomerService_GetCustomer_Failed(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)

	mockRepo.EXPECT().ById("1").Return(nil, errs.NewUnexpectedError("Unexpected database error"))

	// Act
	_, err := service.GetCustomer("1")

	// Assert
	if want := errs.NewUnexpectedError("Unexpected database error"); !reflect.DeepEqual(err, want) {
		t.Errorf("Test GetAllCustomer failed:\ngot: %v\nwant: %v\n", err, want)
	}
}

func TestDefaultCustomerService_SaveCustomer(t *testing.T) {
	type fields struct {
		repo domain.CustomerRepository
	}
	type args struct {
		req dto.CustomerRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *dto.CustomerResponse
		want1  *errs.AppError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DefaultCustomerService{
				repo: tt.fields.repo,
			}
			got, got1 := s.SaveCustomer(tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveCustomer() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SaveCustomer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewCustomerService(t *testing.T) {
	type args struct {
		repository domain.CustomerRepository
	}
	tests := []struct {
		name string
		args args
		want DefaultCustomerService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCustomerService(tt.args.repository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCustomerService() = %v, want %v", got, tt.want)
			}
		})
	}
}
