package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"red/errs"
	"time"
)

var now func() time.Time

func init() {
	now = func() time.Time {
		return time.Now()
	}
}

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
		//logger_zap.Error("Error while creating new account")
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}

	return &a, nil
}

func (d AccountRepositoryDB) Update(a Account) (*Account, *errs.AppError) {
	// 將帳戶資料存進DB
	result := d.client.Model(&Account{}).Where("account_id", a.AccountId).Updates(&a)
	err := result.Error
	if err != nil {
		//logger_zap.Error("Error while creating new account")
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}

	return &a, nil
}

func (d AccountRepositoryDB) Delete(accountID string) *errs.AppError {
	deleteDate := time.Now()
	result := d.client.Model(&Account{}).Where("account_id = ?", accountID).Updates(Account{Status: "0", DeleteAt: &deleteDate})
	if err := result.Error; err != nil {
		return errs.NewUnexpectedError("Unexpected database error when soft delete account")
	}

	return nil
}

// ByID 找尋特定id的帳戶資料
func (d AccountRepositoryDB) ByID(accountID uint) (*Account, *errs.AppError) {
	var a Account
	// 在account表格裡預載入交易紀錄的資料，並且讀取特定id的資料
	result := d.client.Where("account_id = ?", accountID).First(&a)
	//result := d.client.Table("accounts").Preload("Transactions").Where("account_id = ?", id).Find(&a)
	if result.RowsAffected == 0 {
		return nil, errs.NewNotFoundError("Account not found")
	}

	if err := result.Error; err != nil {
		//logger_zap.Error("Error while querying accounts table" + err.Error())
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Account not found")
		}
		return nil, errs.NewUnexpectedError("Unexpected database error when find account by account id")
	}
	return &a, nil
}

func (d AccountRepositoryDB) TransactionsByID(accountID uint) ([]Transaction, *errs.AppError) {
	var t []Transaction
	result := d.client.Where("account_id = ?", accountID).Find(&t)
	//if result.RowsAffected == 0 {
	//	return nil, errs.NewNotFoundError("Transactions not found")
	//}

	if err := result.Error; err != nil {
		//logger_zap.Error("Error while querying accounts table" + err.Error())
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("transactions not found")
		}
		return nil, errs.NewUnexpectedError("Unexpected database error when find transactions by account id")
	}

	return t, nil
}

// ByCustomerID 找尋特定id的帳戶資料
func (d AccountRepositoryDB) ByCustomerID(id uint) ([]Account, *errs.AppError) {
	var a []Account

	result := d.client.Where("customer_id = ?", id).Find(&a)

	if err := result.Error; err != nil {
		//logger_zap.Error("Error while querying accounts table" + err.Error())
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
		//logger_zap.Error("Error while starting a new transaction for bank account transaction" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error when save transaction")
	}

	if t.IsWithdrawal() {
		t.NewBalance -= t.Amount
	} else {
		t.NewBalance += t.Amount
	}

	//result := tx.Table("accounts").Where("account_id", t.AccountId).Update("amount", t.NewBalance)
	result := tx.Model(&Account{}).Where("account_id", t.AccountId).
		Updates(Account{Amount: t.NewBalance})
	if err := result.Error; err != nil {
		tx.Rollback()
		//logger_zap.Error("Error while update account amount" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error when save transaction")
	}

	result = tx.Create(&t)
	// 更改帳戶金額失敗則回滾
	if err := result.Error; err != nil {
		tx.Rollback()
		//logger_zap.Error("Error while update account amount" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error when save transaction")
	}
	// 進行提交
	result = tx.Commit()
	if err = result.Error; err != nil {
		tx.Rollback()
		//logger_zap.Error("Error while committing transaction for bank account" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error while saving transaction commit failed ")
	}

	return &t, nil
}
