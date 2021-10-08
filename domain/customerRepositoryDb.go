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

func (d CustomerRepositoryDb) Save(c Customer) (*Customer, *errs.AppError) {
	//fmt.Println(c)
	result := d.client.Create(&c)
	err := result.Error
	if err != nil {
		logger.Error("Error while creating new customer")
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}

	return &c, nil
}

func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {

	var c Customer
	result := d.client.Table("Customers").Preload("Accounts").Where("customer_id = ?", id).Find(&c)
	if result.Error != nil {
		if result.Error == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		} else {
			logger.Error("Error while querying customers table" + result.Error.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}

	//fmt.Println(c)
	return &c, nil
}

func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {
	var err error
	customers := make([]Customer, 0)

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
