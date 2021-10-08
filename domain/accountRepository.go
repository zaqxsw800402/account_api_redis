package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"red/errs"
	"red/logger"
)

type AccountRepositoryDB struct {
	client *gorm.DB
}

func NewAccountRepositoryDb(dbClient *gorm.DB) AccountRepositoryDB {
	return AccountRepositoryDB{dbClient}
}

func (d AccountRepositoryDB) Save(a Account) (*Account, *errs.AppError) {

	result := d.client.Create(&a)
	err := result.Error
	if err != nil {
		logger.Error("Error while creating new account")
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}

	return &a, nil
}

func (d AccountRepositoryDB) FindBy(id string) (*Account, *errs.AppError) {
	var a Account
	result := d.client.Where("account_id = ?", id).First(&a)
	err := result.Error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Account not found")
		} else {
			logger.Error("Error while querying accounts table" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &a, nil
}

func (d AccountRepositoryDB) SaveTransaction(t Transaction) (*Transaction, *errs.AppError) {
	tx := d.client.Begin()
	err := tx.Error
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account transaction" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	result := tx.Create(&t)

	if t.IsWithdrawal() {
		result = tx.Table("accounts").Where("account_id", t.AccountId).Update("amount", gorm.Expr("amount - ?", t.Amount))
	} else {
		result = tx.Table("accounts").Where("account_id", t.AccountId).Update("amount", gorm.Expr("amount + ?", t.Amount))
	}

	err = result.Error
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	result = tx.Commit()
	err = result.Error
	if err != nil {
		tx.Rollback()
		logger.Error("Error while committing transaction for bank account" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	account, appError := d.FindBy(t.AccountId)
	if appError != nil {
		return nil, appError
	}

	t.Amount = account.Amount
	return &t, nil
}
