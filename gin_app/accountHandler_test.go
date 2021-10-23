package gin_app

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"red/Redis"
	"red/domain"
	"red/dto"
	"red/errs"
	service2 "red/mocks/service"
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

	account := dto.NewAccountRequest{
		CustomerId:  3,
		AccountType: "saving",
		Amount:      6000,
	}

	expectedAccount := dto.NewAccountResponse{
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

	account := dto.NewAccountRequest{
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

	body := `{
    "amount" : 6000
}`
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
	///customers/:id/account/:account_id"

	request, _ := http.NewRequest(http.MethodGet, "/customers/3/account/1", nil)

	expectedAccount := domain.Account{
		AccountId:    1,
		CustomerId:   3,
		OpeningDate:  "2021-10-08 14:48:40",
		AccountType:  "saving",
		Amount:       6000,
		Status:       "1",
		Transactions: nil,
	}

	mockAccount.EXPECT().GetAccount(uint(1)).Return(&expectedAccount, nil)
	router.GET("/customers/:id/account/:account_id", ah.getAccount)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := `{"account_id":1,"customer_id":3,"opening_date":"2021-10-08 14:48:40","account_type":"saving","amount":6000,"status":"1","transactions":null}`

	if exp := recorder.Body.String(); exp != response {
		t.Errorf("Failed test while get account, \nexpected: %v\ngot: %v", exp, response)
	}
}

//func TestAccountHandler_getAccount(t *testing.T) {
//	//type fields struct {
//	//	service service.AccountService
//	//	redisDB Redis.Database
//	//}
//	//type args struct {
//	//	c *gin.Context
//	//}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := AccountHandler{
//				service: tt.fields.service,
//				redisDB: tt.fields.redisDB,
//			}
//		})
//	}
//}

//func TestAccountHandler_makeTransaction(t *testing.T) {
//	type fields struct {
//		service service.AccountService
//		redisDB Redis.Database
//	}
//	type args struct {
//		c *gin.Context
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := AccountHandler{
//				service: tt.fields.service,
//				redisDB: tt.fields.redisDB,
//			}
//		})
//	}
//}
