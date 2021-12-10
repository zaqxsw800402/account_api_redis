package service

import (
	"github.com/golang/mock/gomock"
	domain3 "red/cmd/api/domain"
	dto2 "red/cmd/api/dto"
	"red/cmd/api/errs"
	domain2 "red/mocks/domain"
	"reflect"
	"testing"
	"time"
)

func TestNewAccount_Validate_Failed(t *testing.T) {
	//Arrange
	request := dto2.AccountRequest{
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
	req := dto2.AccountRequest{
		CustomerId:  2,
		AccountType: "saving",
		Amount:      6000,
	}
	account := domain3.Account{
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
	req := dto2.AccountRequest{
		CustomerId:  2,
		AccountType: "saving",
		Amount:      6000,
	}
	account := domain3.Account{
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
	account := &domain3.Account{
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
	request := dto2.TransactionRequest{
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

func TestMakeTransaction_DepositSuccess(t *testing.T) {
	//Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockAccountRepository(ctrl)
	service := NewAccountService(mockRepo)

	account := domain3.Account{
		AccountId:   1,
		CustomerId:  1,
		OpeningDate: "",
		AccountType: "saving",
		Amount:      6000,
		Status:      "1",
	}

	req := dto2.TransactionRequest{
		AccountId:       1,
		Amount:          7000,
		TransactionType: "deposit",
	}
	mockRepo.EXPECT().FindBy(req.AccountId).Return(&account, nil)

	transaction := domain3.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	returnTransaction := &domain3.Transaction{
		TransactionId:   1,
		AccountId:       req.AccountId,
		Amount:          13000,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	want := &dto2.TransactionResponse{
		TransactionId:   1,
		AccountId:       req.AccountId,
		Amount:          13000,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	mockRepo.EXPECT().SaveTransaction(transaction).Return(returnTransaction, nil)
	//mockRepo.SaveTransaction(transaction)
	//Act
	got, appError := service.MakeTransaction(req)
	//Assert
	if appError != nil {
		t.Error("failed while test the account amount validation")
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("failed while test the transaction, expected:%v , got:%v", want, got)
	}
}

func TestMakeTransaction_SaveFailed(t *testing.T) {
	//Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockAccountRepository(ctrl)
	service := NewAccountService(mockRepo)

	account := domain3.Account{
		AccountId:   1,
		CustomerId:  1,
		OpeningDate: "",
		AccountType: "saving",
		Amount:      6000,
		Status:      "1",
	}

	req := dto2.TransactionRequest{
		AccountId:       1,
		Amount:          7000,
		TransactionType: "deposit",
	}
	mockRepo.EXPECT().FindBy(req.AccountId).Return(&account, nil)

	transaction := domain3.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	mockRepo.EXPECT().SaveTransaction(transaction).Return(nil, errs.NewUnexpectedError("Unexpected database error"))
	//mockRepo.SaveTransaction(transaction)
	//Act
	_, err := service.MakeTransaction(req)
	//Assert
	if want := errs.NewUnexpectedError("Unexpected database error"); !reflect.DeepEqual(err, want) {
		t.Errorf("failed while check the account amount validation. \ngot: %v\nwant: %v", err, want)
	}
}

func TestMakeTransaction_WithdrawalSuccess(t *testing.T) {
	//Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockAccountRepository(ctrl)
	service := NewAccountService(mockRepo)

	account := domain3.Account{
		AccountId:   1,
		CustomerId:  1,
		OpeningDate: "",
		AccountType: "saving",
		Amount:      16000,
		Status:      "1",
	}

	req := dto2.TransactionRequest{
		AccountId:       1,
		Amount:          7000,
		TransactionType: "withdrawal",
	}
	mockRepo.EXPECT().FindBy(req.AccountId).Return(&account, nil)

	transaction := domain3.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	returnTransaction := &domain3.Transaction{
		TransactionId:   1,
		AccountId:       req.AccountId,
		Amount:          9000,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	want := &dto2.TransactionResponse{
		TransactionId:   1,
		AccountId:       req.AccountId,
		Amount:          9000,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	mockRepo.EXPECT().SaveTransaction(transaction).Return(returnTransaction, nil)
	//mockRepo.SaveTransaction(transaction)
	//Act
	got, appError := service.MakeTransaction(req)
	//Assert
	if appError != nil {
		t.Error("failed while test the account amount validation")
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("failed while test the transaction, expected:%v , got:%v", want, got)
	}
}

func TestMakeTransaction_WithdrawalFailed(t *testing.T) {
	//Arrange
	ctrl := gomock.NewController(t)
	mockRepo := domain2.NewMockAccountRepository(ctrl)
	service := NewAccountService(mockRepo)

	account := domain3.Account{
		AccountId:   1,
		CustomerId:  1,
		OpeningDate: "",
		AccountType: "saving",
		Amount:      16000,
		Status:      "1",
	}

	req := dto2.TransactionRequest{
		AccountId:       1,
		Amount:          17000,
		TransactionType: "withdrawal",
	}
	mockRepo.EXPECT().FindBy(req.AccountId).Return(&account, nil)

	//Act
	_, err := service.MakeTransaction(req)
	//Assert
	if want := errs.NewValidationError("Insufficient balance in the account"); !reflect.DeepEqual(err, want) {
		t.Errorf("failed while check the account amount validation. \ngot: %v\nwant: %v", err, want)
	}

}
