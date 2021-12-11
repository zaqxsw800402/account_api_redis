package service

import (
	"red/cmd/api/domain"
	"red/cmd/api/dto"
	"red/cmd/api/errs"
	"time"
)

const dbTSLayout = "2006-01-02 15:04:05"

//go:generate mockgen -destination=../mocks/service/mockAccountService.go -package=service red/service AccountService
type AccountService interface {
	NewAccount(request dto.AccountRequest) (*dto.AccountResponse, *errs.AppError)
	GetAccount(customerID uint, accountId uint) (*dto.AccountResponse, *errs.AppError)
	GetALlTransactions(customerID uint, accountId uint) ([]dto.TransactionResponse, *errs.AppError)
	GetAllAccount(customerId uint) ([]dto.AccountResponse, *errs.AppError)
	MakeTransaction(request dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repo domain.AccountRepository
}

// NewAccountService 建立結構體
func NewAccountService(repo domain.AccountRepository) DefaultAccountService {
	return DefaultAccountService{repo}
}

// GetAccount 藉由repository讀取特定的帳戶資料
func (s DefaultAccountService) GetAccount(customerID uint, accountId uint) (*dto.AccountResponse, *errs.AppError) {
	// 讀取特定帳戶ID的資料
	account, err := s.repo.ByID(customerID, accountId)
	if err != nil {
		return nil, err
	}
	response := account.ToDto()
	return &response, nil
}

func (s DefaultAccountService) GetALlTransactions(customerID uint, accountId uint) ([]dto.TransactionResponse, *errs.AppError) {
	transactions, err := s.repo.TransactionsByID(accountId)
	if err != nil {
		return nil, err
	}
	response := make([]dto.TransactionResponse, 0)
	for _, transaction := range transactions {
		response = append(response, transaction.ToDto())
	}
	return response, nil
}

func (s DefaultAccountService) GetAllAccount(customerID uint) ([]dto.AccountResponse, *errs.AppError) {
	// 讀取特定帳戶ID的資料
	accounts, err := s.repo.ByCustomerID(customerID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.AccountResponse, 0)
	for _, account := range accounts {
		response = append(response, account.ToDto())
	}
	return response, nil
}

// NewAccount 建立新帳戶
func (s DefaultAccountService) NewAccount(req dto.AccountRequest) (*dto.AccountResponse, *errs.AppError) {
	// 檢查body裡的資料是否有效
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	// 存進結構體，並存入db
	a := domain.Account{
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format(dbTSLayout),
		AccountType: req.AccountType,
		Amount:      float64(req.Amount),
		Status:      "1",
	}
	newAccount, err := s.repo.Save(a)
	if err != nil {
		return nil, err
	}

	// 轉換成回傳的json格式
	response := newAccount.ToNewAccountResponseDto()
	return &response, nil
}

// MakeTransaction 判斷是否能交易，以及儲存交易紀錄
func (s DefaultAccountService) MakeTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError) {
	// 判斷內容是否有效
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	t := domain.Transaction{
		AccountId:       uint(req.AccountId),
		Amount:          float64(req.Amount),
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}

	customerID := req.CustomerId
	// 取出帳戶金額
	account, err := s.repo.ByID(uint(customerID), t.AccountId)
	if err != nil {
		return nil, err
	}
	// 判斷金額
	if req.IsTransactionTypeWithdrawal() {
		if !account.CanWithdraw(t.Amount) {
			return nil, errs.NewValidationError("Insufficient balance in the account")
		}
	}

	// 存入交易紀錄
	t.NewBalance = account.Amount
	transaction, appError := s.repo.SaveTransaction(account.CustomerId, t)
	if appError != nil {
		return nil, appError
	}

	// 轉換成回傳的json格式
	response := transaction.ToDto()

	return &response, nil
}
