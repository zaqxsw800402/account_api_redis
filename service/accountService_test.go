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

func TestNewAccount_Validate_Failed(t *testing.T) {
	//Arrange
	request := dto.NewAccountRequest{
		CustomerId:  2,
		AccountType: "saving",
		Amount:      0,
	}
	service := NewAccountService(nil)

	//Act
	_, appError := service.NewAccount(request)

	//Assert
	if appError == nil {
		t.Error("failed while test the new account validation")
	}
}

func TestNewAccount_Create_Fail(t *testing.T) {
	//Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockAccountRepository(ctrl)
	service := NewAccountService(mockRepo)
	req := dto.NewAccountRequest{
		CustomerId:  2,
		AccountType: "saving",
		Amount:      6000,
	}
	account := domain.Account{
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format(dbTSLayout),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}
	mockRepo.EXPECT().Save(account).Return(nil, errs.NewUnexpectedError("Unexpected database error"))
	//Act
	_, appError := service.NewAccount(req)

	//Assert
	if appError == nil {
		t.Errorf("Test failed while validating error for new account")
	}
}

func TestNewAccount_Create_Success(t *testing.T) {
	//Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockAccountRepository(ctrl)
	service := NewAccountService(mockRepo)
	req := dto.NewAccountRequest{
		CustomerId:  2,
		AccountType: "saving",
		Amount:      6000,
	}
	account := domain.Account{
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format(dbTSLayout),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}
	accountWithId := account
	accountWithId.AccountId = 201
	mockRepo.EXPECT().Save(account).Return(&accountWithId, nil)
	//Act
	newAccount, appError := service.NewAccount(req)

	//Assert
	if appError != nil {
		t.Errorf("Test failed while validating error for new account")
	}
	if newAccount.AccountId != accountWithId.AccountId {
		t.Error("Test failed while matching new account id")
	}
}

func TestGetAccount_Success(t *testing.T) {
	//Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockAccountRepository(ctrl)
	service := NewAccountService(mockRepo)
	account := &domain.Account{
		AccountId:   1,
		CustomerId:  201,
		OpeningDate: "",
		AccountType: "",
		Amount:      6000,
		Status:      "",
	}
	mockRepo.EXPECT().FindBy(uint(1)).Return(account, nil)
	//Act
	newAccount, appError := service.GetAccount(1)
	//Assert
	if appError != nil {
		t.Error("Test failed while validating error for get account with id")
	}
	if newAccount != account {
		t.Error("Test failed while get account with id")
	}
}

func TestGetAccount_Failed(t *testing.T) {
	//Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockAccountRepository(ctrl)
	service := NewAccountService(mockRepo)

	mockRepo.EXPECT().FindBy(uint(1)).Return(nil, errs.NewNotFoundError("Account not found"))
	//Act
	_, appError := service.GetAccount(1)
	//Assert
	if appError == nil {
		t.Error("Test failed while validating error for get account with id")
	}

}

func TestMakeTransaction_Validate_Failed(t *testing.T) {
	//Arrange
	request := dto.TransactionRequest{
		AccountId:       1,
		Amount:          6000,
		TransactionType: "d",
	}
	service := NewAccountService(nil)

	//Act
	_, appError := service.MakeTransaction(request)

	//Assert
	if appError == nil {
		t.Error("failed while test the new account validation")
	}
}

func TestMakeTransaction_Type_Failed(t *testing.T) {
	//Arrange
	req := dto.TransactionRequest{
		AccountId:       1,
		Amount:          7000,
		TransactionType: "deposit",
	}
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockAccountRepository(ctrl)
	service := NewAccountService(mockRepo)
	transaction := domain.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	saveTransaction := &domain.Transaction{
		TransactionId:   1,
		AccountId:       req.AccountId,
		Amount:          12000,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	dtoTransaction := &dto.TransactionResponse{
		TransactionId:   1,
		AccountId:       req.AccountId,
		Amount:          12000,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	mockRepo.EXPECT().SaveTransaction(transaction).Return(saveTransaction, nil)
	//Act
	newTransaction, appError := service.MakeTransaction(req)
	//Assert
	if appError != nil {
		t.Error("failed while test the account amount validation")
	}

	if !reflect.DeepEqual(dtoTransaction, newTransaction) {
		t.Errorf("failed while test the transaction, expected:%v , got:%v", dtoTransaction, newTransaction)
	}
}
