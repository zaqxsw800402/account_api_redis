package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"red/dto"
	"red/errs"
	service2 "red/mocks/service"
	"red/redis"
	"testing"
)

var router *gin.Engine
var app application
var mockCustomer *service2.MockCustomerService
var mockAccount *service2.MockAccountService

func setUp(t *testing.T) {
	router = gin.Default()
	// mock service
	ctrl := gomock.NewController(t)
	mockCustomer = service2.NewMockCustomerService(ctrl)
	mockAccount = service2.NewMockAccountService(ctrl)

	client, _ := redismock.NewClientMock()
	db := redis.New(client)
	//	db := redis.New(client)

	ch := CustomerHandler{
		mockCustomer,
	}
	ah := AccountHandler{
		mockAccount,
	}

	app = application{ch: ch, ah: ah, redis: db}

}

func TestCustomerHandlers_GetAllCustomer_Success(t *testing.T) {
	setUp(t)

	dummyCustomer := []dto.CustomerResponse{
		{1, "Ivy", "Taiwan", "2000-01-01", "1"},
		{1, "Lily", "Taiwan", "2000-01-01", "1"},
	}
	mockCustomer.EXPECT().GetAllCustomers(0).Return(dummyCustomer, nil)

	router.GET("/all-customers", app.getAllCustomers)

	request, _ := http.NewRequest(http.MethodGet, "/all-customers", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code")
	}

}

func TestCustomerHandlers_GetAllCustomer_Failed_code500(t *testing.T) {
	setUp(t)

	mockCustomer.EXPECT().GetAllCustomers(0).Return(nil, errs.NewUnexpectedError("Unexpected database error"))

	router.GET("/all-customers", app.getAllCustomers)

	request, _ := http.NewRequest(http.MethodGet, "/all-customers", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	//fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code")
	}
}

func TestCustomerHandlers_NewCustomer_Success(t *testing.T) {
	// Arrange
	setUp(t)
	customer := dto.CustomerRequest{
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "1978-12-15",
	}

	expectedCustomer := dto.CustomerResponse{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "1978-12-15",
		Status:      "active",
	}
	mockCustomer.EXPECT().SaveCustomer(0, customer).Return(&expectedCustomer, nil)
	router.POST("/all-customers/0", app.newCustomer)

	body := `{
		"name": "Ivy",
		"city": "Taiwan",
		"date_of_birth": "1978-12-15"
}`
	b := bytes.NewBufferString(body)

	// Act
	request, _ := http.NewRequest(http.MethodPost, "/all-customers/0", b)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	//response := `{"customer_id":1,"name":"Ivy","city":"Taiwan","zipcode":"110075","date_of_birth":"1978-12-15","status":"active","accounts":null}`
	response := `{"error":false,"message":""}`
	// Assert
	//fmt.Println(recorder.Body.String())
	if exp := recorder.Body.String(); exp != response {
		t.Errorf("Failed test while create new customer, \nexpected: %v\ngot: %v", exp, response)
	}

}

func TestCustomerHandlers_NewCustomer_FailedSave(t *testing.T) {
	// Arrange
	setUp(t)
	customer := dto.CustomerRequest{
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "1978-12-15",
	}

	mockCustomer.EXPECT().SaveCustomer(0, customer).Return(nil, errs.NewUnexpectedError("Unexpected error from database"))
	router.POST("/all-customers/0", app.newCustomer)

	body := `{
		"name": "Ivy",
		"city": "Taiwan",
		"date_of_birth": "1978-12-15"
}`

	b := bytes.NewBufferString(body)

	// Act
	request, _ := http.NewRequest(http.MethodPost, "/all-customers/0", b)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := `{"error":true,"message":"Unexpected error from database"}`

	// Assert
	fmt.Println(recorder.Body.String())
	if exp := recorder.Body.String(); exp != response {
		t.Errorf("Failed test while create save customer, \nexpected: %v\ngot: %v", exp, response)
	}

}

func TestCustomerHandlers_NewCustomer_Failed(t *testing.T) {
	// Arrange
	setUp(t)
	router.POST("/all-customers/0", app.newCustomer)

	body := `{
		"name": "Ivy",
		"date_of_birth": "1978-12-15"
}`
	b := bytes.NewBufferString(body)

	// Act
	request, _ := http.NewRequest(http.MethodPost, "/all-customers/0", b)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	//fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusBadRequest {
		t.Error("Failed while testing the status code")
	}
}

func TestCustomerHandlers_GetCustomer_Success(t *testing.T) {
	// Arrange
	setUp(t)

	expectedCustomer := dto.CustomerResponse{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "1978-12-15",
		Status:      "active",
	}
	mockCustomer.EXPECT().GetCustomer("").Return(&expectedCustomer, nil)

	router.GET("/api/admin/all-customers/1", app.getCustomer)

	response := `{"customer_id":1,"name":"Ivy","city":"Taiwan","date_of_birth":"1978-12-15","status":"active"}`

	// Act
	request, _ := http.NewRequest(http.MethodGet, "/api/admin/all-customers/1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	fmt.Println(recorder.Body.String())
	if exp := recorder.Body.String(); exp != response {
		t.Errorf("Failed test while create new customer, \nexpected: %v\ngot: %v", exp, response)
	}

}

func TestCustomerHandlers_GetCustomer_BadRequest(t *testing.T) {
	// Arrange
	setUp(t)

	mockCustomer.EXPECT().GetCustomer("30").Return(nil, errs.NewUnexpectedError("Unexpected database error"))
	router.GET("/customers/:id", app.getCustomer)

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
