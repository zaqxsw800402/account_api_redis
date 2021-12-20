package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"red/dto"
	"red/errs"
	"reflect"
	"testing"
)

func TestAccountHandler_NewAccount_Success(t *testing.T) {
	//Arrange
	setUp(t)

	body := `{
	"customer_id": 1,
	"account_type": "saving",
	"amount" : 6000
}`
	b := bytes.NewBufferString(body)

	expectedCustomer := dto.CustomerResponse{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "1978-12-15",
		Status:      "active",
	}
	mockCustomer.EXPECT().GetCustomer("1").Return(&expectedCustomer, nil)

	account := dto.AccountRequest{
		CustomerId:  1,
		AccountType: "saving",
		Amount:      6000,
	}

	expectedAccount := dto.AccountResponse{
		AccountId:   1,
		CustomerId:  1,
		AccountType: "saving",
		Amount:      6000,
		Status:      "active",
	}

	mockAccount.EXPECT().NewAccount(account).Return(&expectedAccount, nil)

	request, err := http.NewRequest(http.MethodPost, "/api/admin/all-customers/accounts/0", b)
	if err != nil {
		fmt.Println(err)
	}

	router.POST("/api/admin/all-customers/accounts/0", app.newAccount)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := `{"error":false,"message":""}`

	if exp := recorder.Body.String(); exp != response {
		t.Errorf("Failed test while create new account, \nexpected: %v\ngot: %v", exp, response)
	}
}

func TestAccountHandler_NewAccount_Failed_inactive(t *testing.T) {
	//Arrange
	setUp(t)

	body := `{
	"customer_id": 1,
	"account_type": "saving",
	"amount" : 6000
}`
	b := bytes.NewBufferString(body)
	request, _ := http.NewRequest(http.MethodPost, "/api/admin/all-customers/accounts", b)

	expectedCustomer := dto.CustomerResponse{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		DateOfBirth: "1978-12-15",
		Status:      "inactive",
	}
	mockCustomer.EXPECT().GetCustomer("1").Return(&expectedCustomer, nil)

	router.POST("/api/admin/all-customers/accounts", app.newAccount)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := `{"error":true,"message":"this customer is inactive"}`

	if exp := recorder.Body.String(); exp != response {
		t.Errorf("Failed test while create new account, \nexpected: %v\ngot: %v", exp, response)
	}
}

func TestAccountHandler_NewAccount_Failed_NotFound(t *testing.T) {
	//Arrange
	setUp(t)

	body := `{
	"customer_id": 1,
	"account_type": "saving",
	"amount" : 6000
}`
	b := bytes.NewBufferString(body)
	request, _ := http.NewRequest(http.MethodPost, "/api/admin/all-customers/accounts", b)

	mockCustomer.EXPECT().GetCustomer("1").Return(nil, errs.NewNotFoundError("account not found"))

	router.POST("/api/admin/all-customers/accounts", app.newAccount)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := `{"error":true,"message":"account not found"}`

	if exp := recorder.Body.String(); exp != response {
		t.Errorf("Failed test while create new account, \nexpected: %v\ngot: %v", exp, response)
	}
}

func TestAccountHandler_NewAccount_BadRequest(t *testing.T) {
	//Arrange
	setUp(t)

	body := `{"amount" : 6000}`
	b := bytes.NewBufferString(body)
	request, _ := http.NewRequest(http.MethodPost, "/customers/3/account", b)

	router.POST("/customers/:id/account", app.newAccount)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Error("Failed while testing the status code")
	}
}

func TestAccountHandler_makeTransaction_Success(t *testing.T) {
	// Arrange
	setUp(t)

	router.POST("/api/admin/withdrawal", app.makeTransaction)

	body := `{
	"customer_id":1,
	"account_id":1,
	"target_account_id":1,
	"transaction_type":"deposit",
	"amount" : 6000
}`
	b := bytes.NewBufferString(body)

	req := dto.TransactionRequest{
		AccountId:       1,
		Amount:          6000,
		TransactionType: "deposit",
		TransactionDate: "",
		CustomerId:      1,
		TargetAccountId: 1,
	}

	res := dto.TransactionResponse{
		TransactionId:   1,
		AccountId:       1,
		NewBalance:      7000,
		TransactionType: "deposit",
		TransactionDate: "",
	}
	mockAccount.EXPECT().MakeTransaction(req).Return(&res, nil)

	// Act
	request, _ := http.NewRequest(http.MethodPost, "/api/admin/withdrawal", b)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	exp := `{"error":false,"message":""}`
	if got := recorder.Body.String(); !reflect.DeepEqual(got, exp) {
		t.Errorf("Failed test while get account, \nexpected: %v\ngot: %v", got, exp)
	}
}

func TestAccountHandler_makeTransaction_BadRequest(t *testing.T) {
	// Arrange
	setUp(t)

	router.POST("/api/admin/withdrawal", app.makeTransaction)

	body := `{
	"customer_id":1,
	"account_id":1,
	"target_account_id":1,
	"transaction_type":"dep",
	"amount" : 6000
}`
	b := bytes.NewBufferString(body)

	req := dto.TransactionRequest{
		AccountId:       1,
		Amount:          6000,
		TransactionType: "dep",
		CustomerId:      1,
		TargetAccountId: 1,
	}

	mockAccount.EXPECT().MakeTransaction(req).Return(nil, errs.NewValidationError("transaction type can only be deposit or withdrawal"))

	// Act
	request, _ := http.NewRequest(http.MethodPost, "/api/admin/withdrawal", b)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	exp := `{"error":true,"message":"transaction type can only be deposit or withdrawal"}`
	if got := recorder.Body.String(); !reflect.DeepEqual(got, exp) {
		t.Errorf("Failed test while get account, \nexpected: %v\ngot: %v", got, exp)
	}
}
