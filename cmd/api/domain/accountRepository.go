package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"red/cmd/api/errs"
)

type AccountRepositoryDB struct {
	client *gorm.DB
}

func NewAccountRepositoryDb(dbClient *gorm.DB) AccountRepositoryDB {
	return AccountRepositoryDB{dbClient}
}

// Save 將帳戶資料存進DB
func (d AccountRepositoryDB) Save(a Account) (*Account, *errs.AppError) {
	// 將帳戶資料存進DB
	result := d.client.Create(&a)
	err := result.Error
	if err != nil {
		//logger.Error("Error while creating new account")
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}

	return &a, nil
}

// ByID 找尋特定id的帳戶資料
func (d AccountRepositoryDB) ByID(id uint) (*Account, *errs.AppError) {
	var a Account
	// 在account表格裡預載入交易紀錄的資料，並且讀取特定id的資料
	result := d.client.Table("Accounts").Preload("Transactions").Where("account_id = ?", id).Find(&a)
	if err := result.Error; err != nil {
		//logger.Error("Error while querying accounts table" + err.Error())
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Account not found")
		}
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	return &a, nil
}

// ByCustomerID 找尋特定id的帳戶資料
func (d AccountRepositoryDB) ByCustomerID(id uint) ([]Account, *errs.AppError) {
	var a []Account

	result := d.client.Where("customer_id = ?", id).Find(&a)

	if err := result.Error; err != nil {
		//logger.Error("Error while querying accounts table" + err.Error())
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Account not found with customer_id")
		}
		return nil, errs.NewUnexpectedError("Unexpected database error when find account by using customer_id")
	}
	return a, nil
}

// SaveTransaction 紀錄交易資訊
func (d AccountRepositoryDB) SaveTransaction(t Transaction) (*Transaction, *errs.AppError) {
	// 設定db存入的起始點，用來回滾用
	tx := d.client.Begin()
	err := tx.Error
	if err != nil {
		//logger.Error("Error while starting a new transaction for bank account transaction" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	// 存入交易資訊
	result := tx.Create(&t)
	if err := result.Error; err != nil {
		tx.Rollback()
		//logger.Error("Error while save transaction" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// 更新交易後的帳戶金額
	if t.IsWithdrawal() {
		result = tx.Table("accounts").Where("account_id", t.AccountId).Update("amount", gorm.Expr("amount - ?", t.Amount))
	} else {
		result = tx.Table("accounts").Where("account_id", t.AccountId).Update("amount", gorm.Expr("amount + ?", t.Amount))
	}

	// 更改帳戶金額失敗則回滾
	if err := result.Error; err != nil {
		tx.Rollback()
		//logger.Error("Error while update account amount" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error while updating account amount failed")
	}
	// 進行提交
	result = tx.Commit()
	if err = result.Error; err != nil {
		tx.Rollback()
		//logger.Error("Error while committing transaction for bank account" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error while saving transaction commit failed ")
	}

	// 查詢帳戶金額
	account, appError := d.ByID(t.AccountId)
	if appError != nil {
		return nil, appError
	}
	t.Amount = account.Amount

	return &t, nil
}
