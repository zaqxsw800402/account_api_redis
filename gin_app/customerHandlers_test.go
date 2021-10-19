package gin_app

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"red/Redis"
	"red/dto"
	"red/errs"
	service2 "red/mocks/service"
	"testing"
)

var router *gin.Engine
var ch CustomerHandlers
var mockService *service2.MockCustomerService

//var mock *redismock.ClientMock

func setUp(t *testing.T) {
	router = gin.Default()
	// mock service
	ctrl := gomock.NewController(t)
	mockService = service2.NewMockCustomerService(ctrl)
	// mock redis
	client, _ := redismock.NewClientMock()
	db := Redis.NewRedisDb(client)

	ch = CustomerHandlers{
		mockService,
		db,
	}
}

func TestCustomerHandlers_getAllCustomer_Success(t *testing.T) {
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

func TestCustomerHandlers_getAllCustomer_Failed_code500(t *testing.T) {
	setUp(t)

	mockService.EXPECT().GetAllCustomer("").Return(nil, errs.NewUnexpectedError("Unexpected database error"))

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

func TestCustomerHandlers_newCustomers_Success(t *testing.T) {
	// Arrange
	setUp(t)
	customer := dto.CustomerRequest{
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "110075",
		DateOfBirth: "1978-12-15",
	}

	expectedCustomer := dto.CustomerResponse{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "110075",
		DateOfBirth: "1978-12-15",
		Status:      "active",
		Accounts:    nil,
	}
	mockService.EXPECT().SaveCustomer(customer).Return(&expectedCustomer, nil)
	router.POST("/customers", ch.newCustomers)

	body := `{
		"name": "Ivy",
		"city": "Taiwan",
		"zipcode": "110075",
		"date_of_birth": "1978-12-15"
}`

	response := `{"customer_id":1,"name":"Ivy","city":"Taiwan","zipcode":"110075","date_of_birth":"1978-12-15","status":"active","accounts":null}`
	b := bytes.NewBufferString(body)

	// Act
	request, _ := http.NewRequest(http.MethodPost, "/customers", b)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	fmt.Println(recorder.Body.String())
	if exp := recorder.Body.String(); exp != response {
		t.Errorf("Failed test while create new customer, \nexpected: %v\ngot: %v", exp, response)
	}

}

func TestCustomerHandlers_newCustomers_Failed(t *testing.T) {
	// Arrange
	setUp(t)
	router.POST("/customers", ch.newCustomers)

	body := `{
		"name": "Ivy",
		"city": "Taiwan",
		"date_of_birth": "1978-12-15"
}`
	b := bytes.NewBufferString(body)

	// Act
	request, _ := http.NewRequest(http.MethodPost, "/customers", b)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	//fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusBadRequest {
		t.Error("Failed while testing the status code")
	}
}

func TestCustomerHandlers_getCustomer_Success(t *testing.T) {
	// Arrange
	setUp(t)

	expectedCustomer := dto.CustomerResponse{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "110075",
		DateOfBirth: "1978-12-15",
		Status:      "active",
		Accounts:    nil,
	}
	mockService.EXPECT().GetCustomer("").Return(&expectedCustomer, nil)
	router.GET("/customers/1", ch.getCustomer)

	response := `{"customer_id":1,"name":"Ivy","city":"Taiwan","zipcode":"110075","date_of_birth":"1978-12-15","status":"active","accounts":null}`

	// Act
	request, _ := http.NewRequest(http.MethodGet, "/customers/1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	fmt.Println(recorder.Body.String())
	if exp := recorder.Body.String(); exp != response {
		t.Errorf("Failed test while create new customer, \nexpected: %v\ngot: %v", exp, response)
	}

}

func TestCustomerHandlers_getCustomer_Failed(t *testing.T) {
	// Arrange
	setUp(t)

	mockService.EXPECT().GetCustomer("30").Return(nil, errs.NewUnexpectedError("Unexpected database error"))
	router.GET("/customers/:id", ch.getCustomer)

	// Act
	request, _ := http.NewRequest(http.MethodGet, "/customers/30", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Test failed while validation status code")
	}

}
