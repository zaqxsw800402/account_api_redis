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

func TestDefaultCustomerService_SaveCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)

	customer := domain.Customer{
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "23",
		DateOfBirth: "2006-01-02",
	}
	want := &domain.Customer{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "23",
		DateOfBirth: "2006-01-02",
		Status:      "active",
	}

	mockRepo.EXPECT().Save(customer).Return(want, nil)

	// Act
	req := dto.CustomerRequest{
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "23",
		DateOfBirth: "2006-01-02",
	}
	got, _ := service.SaveCustomer(req)
	if !reflect.DeepEqual(got.Name, want.Name) {
		t.Errorf("Test SaveCustomer: \ngot: %v\nwant %v", got, want)
	}

}

func TestDefaultCustomerService_SaveCustomer_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockCustomerRepository(ctrl)
	service := NewCustomerService(mockRepo)

	customer := domain.Customer{
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "23",
		DateOfBirth: "2006-01-02",
	}

	mockRepo.EXPECT().Save(customer).Return(nil, errs.NewUnexpectedError("Unexpected error from database"))

	// Act
	req := dto.CustomerRequest{
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "23",
		DateOfBirth: "2006-01-02",
	}
	_, err := service.SaveCustomer(req)
	if want := errs.NewUnexpectedError("Unexpected error from database"); !reflect.DeepEqual(err, want) {
		t.Errorf("Test GetAllCustomer failed:\ngot: %v\nwant: %v\n", err, want)
	}

}

func TestDefaultCustomerService_transformStatus(t *testing.T) {
	type fields struct {
		status string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"ActiveSuccess", fields{status: "active"}, "1"},
		{"InactiveSuccess", fields{status: "inactive"}, "0"},
		{"Success", fields{status: ""}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformStatus(tt.fields.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Test transformStatus \ngot: %v, \nwant %v", got, tt.want)
			}
		})
	}
}
