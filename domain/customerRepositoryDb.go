package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"red/errs"
	"red/logger"
)

type CustomerRepositoryDb struct {
	client *gorm.DB
}

func NewCustomerRepositoryDb(dbClient *gorm.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{dbClient}
}

// Save 存入顧客資訊
func (d CustomerRepositoryDb) Save(c Customer) (*Customer, *errs.AppError) {
	// 存入資訊
	result := d.client.Create(&c)
	err := result.Error
	if err != nil {
		logger.Error("Error while creating new customer")
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}

	return &c, nil
}

// ById 找尋特定ID的顧客
func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	// 在customers的表格裡，先預載入account的資料，並讀取特定id的資料
	var c Customer
	result := d.client.Table("Customers").Preload("Accounts").Where("customer_id = ?", id).Find(&c)
	if err := result.Error; err != nil {
		logger.Error("Error while querying customers table" + result.Error.Error())
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		}
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return &c, nil
}

// FindAll 取出所有顧客的資料
func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {
	var err error
	var customers []Customer

	// 根據顧客狀態來找尋資料
	if status == "" {
		result := d.client.Find(&customers)
		err = result.Error
	} else {
		result := d.client.Where("status = ?", status).Find(&customers)
		err = result.Error
	}

	if err != nil {
		logger.Error("Error while querying customers table" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return customers, nil
}
