package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"red/errs"
	"time"
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
		//logger.Error("Error while creating new customer")
		return nil, errs.NewUnexpectedError("Unexpected error from database when create new customer")
	}

	return &c, nil
}

func (d CustomerRepositoryDb) Update(c Customer) (*Customer, *errs.AppError) {
	// 存入資訊
	result := d.client.Model(&Customer{}).Updates(Customer{
		Name:        c.Name,
		City:        c.City,
		Zipcode:     c.Zipcode,
		DateOfBirth: c.DateOfBirth,
	})
	err := result.Error
	if err != nil {
		//logger.Error("Error while creating new customer")
		return nil, errs.NewUnexpectedError("Unexpected error from database when updating customer")
	}

	return &c, nil
}

// ById 找尋特定ID的顧客
func (d CustomerRepositoryDb) ByID(id string) (*Customer, *errs.AppError) {
	// 在customers的表格裡，先預載入account的資料，並讀取特定id的資料
	var c Customer
	result := d.client.Where("customer_id = ?", id).First(&c)
	if err := result.Error; err != nil {
		//logger.Error("Error while querying customers table" + result.Error.Error())
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		}
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return &c, nil
}

// FindAll 取出所有顧客的資料
func (d CustomerRepositoryDb) FindAll(userID int) ([]Customer, *errs.AppError) {
	var customers []Customer

	// 根據顧客狀態來找尋資料
	//if status == "" {
	//	result := d.client.Find(&customers)
	//	err = result.Error
	//} else {
	//	result := d.client.Where("status = ?", status).Find(&customers)
	//	err = result.Error
	//}

	result := d.client.Where(Customer{UserID: uint(userID)}).Find(&customers)

	if err := result.Error; err != nil {
		//logger.Error("Error while querying customers table" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error when find all customer with user id")
	}

	return customers, nil
}

func (d CustomerRepositoryDb) Delete(id string) *errs.AppError {
	//i, err := strconv.ParseUint(id, 10, 10)
	//if err != nil {
	//	return &errs.AppError{
	//		Code:    http.StatusBadRequest,
	//		Message: "could not convert customer id to int",
	//	}
	//}

	deleteDate := time.Now()
	result := d.client.Model(&Customer{}).Where("customer_id = ?", id).Updates(Customer{Status: "0", DeleteAt: &deleteDate})
	if err := result.Error; err != nil {
		return errs.NewUnexpectedError("Unexpected database error when soft delete customer")
	}

	return nil
}
