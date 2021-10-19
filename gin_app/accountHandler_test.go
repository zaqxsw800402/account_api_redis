package gin_app

import (
	"bytes"
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
