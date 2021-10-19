package gin_app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"red/dto"
	"red/errs"
	service2 "red/mocks/service"
	"testing"
)

var router *gin.Engine
var ch CustomerHandlers
var mockService *service2.MockCustomerService

func setUp(t *testing.T) {
	router = gin.Default()
	ctrl := gomock.NewController(t)
	mockService = service2.NewMockCustomerService(ctrl)
	ch = CustomerHandlers{
		mockService,
	}
}

func TestCustomerHandlers_getCustomer(t *testing.T) {
	setUp(t)

	dummyCustomer := []dto.CustomerResponse{
		{1, "Ivy", "Taiwan", "238", "2000-01-01", "1", nil},
		{1, "Lily", "Taiwan", "1111", "2000-01-01", "1", nil},
	}
	mockService.EXPECT().GetAllCustomer("").Return(dummyCustomer, nil)

	router.GET("/customers", ch.getAllCustomers)

	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code")
	}

}

func TestCustomerHandlers_getCustomer_fail_code500(t *testing.T) {
	setUp(t)

	mockService.EXPECT().GetAllCustomer("").Return(nil, errs.NewUnexpectedError("Unexpected database error"))
	ch := CustomerHandlers{
		mockService,
	}

	router.GET("/customers", ch.getAllCustomers)

	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code")
	}
}
