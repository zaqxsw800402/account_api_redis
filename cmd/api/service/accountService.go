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
	MakeTransaction(request dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError)
	GetAccount(accountId uint) (*dto.AccountResponse, *errs.AppError)
	GetAllAccount(customerId uint) ([]dto.AccountResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repo domain.AccountRepository
}

// NewAccountService 建立結構體
func NewAccountService(repo domain.AccountRepository) DefaultAccountService {
	return DefaultAccountService{repo}
}

// GetAccount 藉由repository讀取特定的帳戶資料
func (s DefaultAccountService) GetAccount(accountId uint) (*dto.AccountResponse, *errs.AppError) {
	// 讀取特定帳戶ID的資料
	account, err := s.repo.ByID(accountId)
	if err != nil {
		return nil, err
	}
	response := account.ToDto()
	return &response, nil
}

func (s DefaultAccountService) GetAllAccount(accountId uint) ([]dto.AccountResponse, *errs.AppError) {
	// 讀取特定帳戶ID的資料
	accounts, err := s.repo.ByCustomerID(accountId)
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
		Amount:      req.Amount,
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

	// 取出帳戶金額
	account, err := s.repo.ByID(req.AccountId)
	if err != nil {
		return nil, err
	}
	// 判斷金額
	if req.IsTransactionTypeWithdrawal() {
		if !account.CanWithdraw(req.Amount) {
			return nil, errs.NewValidationError("Insufficient balance in the account")
		}
	}

	t := domain.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(dbTSLayout),
	}
	// 存入交易紀錄
	transaction, appError := s.repo.SaveTransaction(t)
	if appError != nil {
		return nil, appError
	}

	// 轉換成回傳的json格式
	response := transaction.ToDto()

	return &response, nil
}
