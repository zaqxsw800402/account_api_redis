package gin_app

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"red/Redis"
	"red/cmd/api/domain"
	dto2 "red/cmd/api/dto"
	"red/cmd/api/errs"
	service2 "red/mocks/service"
	"reflect"
	"testing"
)

func setUpAccount(t *testing.T) {
	router = gin.Default()
	// mock service
	ctrl := gomock.NewController(t)
	mockAccount = service2.NewMockAccountService(ctrl)
	// mock redis
	client, _ := redismock.NewClientMock()
	db := Redis.NewRedisDb(client)

	ah = AccountHandler{
		mockAccount,
		db,
	}
}

func TestAccountHandler_NewAccount_Success(t *testing.T) {
	//Arrange
	//set body
	setUpAccount(t)

	body := `{
    "account_type": "saving",
    "amount" : 6000
}`
	b := bytes.NewBufferString(body)
	request, _ := http.NewRequest(http.MethodPost, "/customers/3/account", b)

	account := dto2.NewAccountRequest{
		CustomerId:  3,
		AccountType: "saving",
		Amount:      6000,
	}

	expectedAccount := dto2.NewAccountResponse{
		AccountId: 1,
	}

	mockAccount.EXPECT().NewAccount(account).Return(&expectedAccount, nil)
	router.POST("/customers/:id/account", ah.newAccount)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := `{"account_id":1}`

	if exp := recorder.Body.String(); exp != response {
		t.Errorf("Failed test while create new account, \nexpected: %v\ngot: %v", exp, response)
	}
}

func TestAccountHandler_NewAccount_FailedCreate(t *testing.T) {
	//Arrange
	//set body
	setUpAccount(t)

	body := `{
    "account_type": "saving",
    "amount" : 6000
}`
	b := bytes.NewBufferString(body)
	request, _ := http.NewRequest(http.MethodPost, "/customers/3/account", b)

	account := dto2.NewAccountRequest{
		CustomerId:  3,
		AccountType: "saving",
		Amount:      6000,
	}

	mockAccount.EXPECT().NewAccount(account).Return(nil, errs.NewUnexpectedError("Unexpected error from database"))
	router.POST("/customers/:id/account", ah.newAccount)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Test failed while validating status code")
	}

}

func TestAccountHandler_NewAccount_BadRequest(t *testing.T) {
	//Arrange
	//set body
	setUpAccount(t)

	body := `{"amount" : 6000}`
	b := bytes.NewBufferString(body)
	request, _ := http.NewRequest(http.MethodPost, "/customers/3/account", b)

	router.POST("/customers/:id/account", ah.newAccount)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Error("Failed while testing the status code")
	}
}

func TestAccountHandler_GetAccount_Success(t *testing.T) {
	//Arrange
	//set body
	setUpAccount(t)

	expectedAccount := domain.Account{
		AccountId:    1,
		CustomerId:   3,
		OpeningDate:  "2021-10-08 14:48:40",
		AccountType:  "saving",
		Amount:       6000,
		Status:       "1",
		Transactions: nil,
	}

	response := `{"account_id":1,"customer_id":3,"opening_date":"2021-10-08 14:48:40","account_type":"saving","amount":6000,"status":"1","transactions":null}`
	router.GET("/customers/:id/account/:account_id", ah.getAccount)

	request, _ := http.NewRequest(http.MethodGet, "/customers/3/account/1", nil)
	mockAccount.EXPECT().GetAccount(uint(1)).Return(&expectedAccount, nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	//assert
	if got := recorder.Body.String(); !reflect.DeepEqual(got, response) {
		t.Errorf("Failed test while get account, \nexpected: %v\ngot: %v", got, response)

	}
}

func TestAccountHandler_GetAccount_FailedGet(t *testing.T) {
	//Arrange
	setUpAccount(t)

	router.GET("/customers/:id/account/:account_id", ah.getAccount)

	request, _ := http.NewRequest(http.MethodGet, "/customers/3/account/1", nil)
	mockAccount.EXPECT().GetAccount(uint(1)).Return(nil, errs.NewNotFoundError("Account not found"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	//assert
	if got := recorder.Code; !reflect.DeepEqual(got, http.StatusNotFound) {
		t.Errorf("Failed test while get account, \nexpected: %v\ngot: %v", got, http.StatusNotFound)
	}
}

func TestAccountHandler_GetAccount_BadRequest(t *testing.T) {
	//Arrange
	setUpAccount(t)

	response := `strconv.ParseUint: parsing "a": invalid syntax`
	router.GET("/customers/:id/account/:account_id", ah.getAccount)

	request, _ := http.NewRequest(http.MethodGet, "/customers/3/account/a", nil)
	//fmt.Println("err: ",err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	//assert
	if got := recorder.Body.String(); !reflect.DeepEqual(got, response) {
		t.Errorf("Failed test while get account, \nexpected: %v\ngot: %v", got, response)

	}
}

func TestAccountHandler_makeTransaction_Success(t *testing.T) {
	// Arrange
	setUpAccount(t)

	router.POST("/customers/:id/account/:account_id", ah.makeTransaction)

	body := `{
    "transaction_type":"deposit",
    "amount" : 6000
}`
	b := bytes.NewBufferString(body)

	req := dto2.TransactionRequest{
		AccountId:       1,
		Amount:          6000,
		TransactionType: "deposit",
		TransactionDate: "",
		CustomerId:      "3",
	}

	res := dto2.TransactionResponse{
		TransactionId:   1,
		AccountId:       1,
		Amount:          7000,
		TransactionType: "deposit",
		TransactionDate: "",
	}
	mockAccount.EXPECT().MakeTransaction(req).Return(&res, nil)

	// Act
	request, _ := http.NewRequest(http.MethodPost, "/customers/3/account/1", b)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	exp := `{"transaction_id":1,"account_id":1,"new_balance":7000,"transaction_type":"deposit","transaction_date":""}`
	if got := recorder.Body.String(); !reflect.DeepEqual(got, exp) {
		t.Errorf("Failed test while get account, \nexpected: %v\ngot: %v", got, exp)
	}
}

func TestAccountHandler_makeTransaction_BadRequest(t *testing.T) {
	// Arrange
	setUpAccount(t)

	router.POST("/customers/:id/account/:account_id", ah.makeTransaction)

	body := `{
    "transaction_type":"dep",
    "amount" : 6000
}`
	b := bytes.NewBufferString(body)

	req := dto2.TransactionRequest{
		AccountId:       1,
		Amount:          6000,
		TransactionType: "dep",
		TransactionDate: "",
		CustomerId:      "3",
	}

	mockAccount.EXPECT().MakeTransaction(req).Return(nil, errs.NewValidationError("transaction type can only be deposit or withdrawal"))

	// Act
	request, _ := http.NewRequest(http.MethodPost, "/customers/3/account/1", b)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	exp := `{"Message":"transaction type can only be deposit or withdrawal"}`
	if got := recorder.Body.String(); !reflect.DeepEqual(got, exp) {
		t.Errorf("Failed test while get account, \nexpected: %v\ngot: %v", got, exp)
	}
}
