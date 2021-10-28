package domain

import (
	"github.com/DATA-DOG/go-sqlmock"
	"red/errs"
	"reflect"
	"regexp"
	"testing"
)

func TestCustomerRepositoryDb_Save(t *testing.T) {
	// Arrange
	client := setup()
	db := NewCustomerRepositoryDb(client)

	mock.ExpectBegin()

	mock.ExpectExec(
		// INSERT INTO `customers` (`name`,`city`,`zipcode`,`date_of_birth`,`status`,`customer_id`) VALUES (?,?,?,?,?,?)
		regexp.QuoteMeta("INSERT INTO `customers` (`name`,`city`,`zipcode`,`date_of_birth`,`status`,`customer_id`) VALUES (?,?,?,?,?,?)")).
		WithArgs("Ivy", "Taiwan", "23", "2012-10-18", "1", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// Act
	customer := Customer{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "23",
		DateOfBirth: "2012-10-18",
		Status:      "1",
	}
	_, err := db.Save(customer)
	// Assert
	if err != nil {
		t.Errorf("test save customer failed: %v", err)
	}
}

func TestCustomerRepositoryDb_Failed(t *testing.T) {
	// Arrange
	client := setup()
	db := NewCustomerRepositoryDb(client)

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta("INSERT INTO `customers` (`name`,`city`,`zipcode`,`date_of_birth`,`status`,`customer_id`) VALUES (?,?,?,?,?,?)")).
		WithArgs("Ivy", "Taiwan", "23", "2012-10-18", 1).
		WillReturnError(errs.NewUnexpectedError("Unexpected error from database"))

	mock.ExpectCommit()

	// Act
	customer := Customer{
		Id:          1,
		Name:        "Ivy",
		City:        "Taiwan",
		Zipcode:     "23",
		DateOfBirth: "2012-10-18",
		Status:      "1",
	}
	_, err := db.Save(customer)

	// Assert
	if want := errs.NewUnexpectedError("Unexpected error from database"); !reflect.DeepEqual(err, want) {
		t.Errorf("validate save customer error failed:\ngot: %v\n want: %v\n", err, want)
	}
}

func TestCustomerRepositoryDb_ById_Success(t *testing.T) {
	// Arrange
	client := setup()
	db := NewCustomerRepositoryDb(client)

	rows := sqlmock.NewRows([]string{`name`, `city`, `zipcode`, `date_of_birth`, `status`, `customer_id`}).
		AddRow("Ivy", "Taiwan", "23", "2012-10-18", "1", 1)

	mock.ExpectQuery(
		"^SELECT \\* FROM `Customers` WHERE customer_id = \\?").
		WithArgs("1").
		WillReturnRows(rows)

	accountRows := sqlmock.NewRows([]string{"customer_id", "opening_date", "account_type", "amount", "status", "account_id"}).
		AddRow(1, "2012-10-18", "saving", 6000, 1, 1)

	mock.ExpectQuery(
		"^SELECT \\* FROM `accounts` WHERE `accounts`.`customer_id` = \\?").
		WillReturnRows(accountRows)

	// Act
	_, err := db.ById("1")

	// Assert
	if err != nil {
		t.Errorf("test find customer id failed: %v", err)
	}
}

func TestCustomerRepositoryDb_ById_Failed(t *testing.T) {
	// Arrange
	client := setup()
	db := NewCustomerRepositoryDb(client)

	mock.ExpectQuery(
		"^SELECT \\* FROM `Customers` WHERE customer_id = \\?").
		WithArgs("1").
		//WillReturnRows(rows)
		WillReturnError(errs.NewUnexpectedError("Unexpected database error"))

	// Act
	_, err := db.ById("1")

	// Assert
	if want := errs.NewUnexpectedError("Unexpected database error"); !reflect.DeepEqual(err, want) {
		t.Errorf("validate save customer error failed:\ngot: %v\n want: %v\n", err, want)
	}
}

func TestCustomerRepositoryDb_FindAll_Success(t *testing.T) {
	// Arrange
	client := setup()
	db := NewCustomerRepositoryDb(client)

	rows := sqlmock.NewRows([]string{`name`, `city`, `zipcode`, `date_of_birth`, `status`, `customer_id`}).
		AddRow("Ivy", "Taiwan", "23", "2012-10-18", "1", 1).
		AddRow("Lily", "Taiwan", "23", "2012-10-18", "1", 1)

	mock.ExpectQuery(
		"^SELECT \\* FROM `customers`").
		WillReturnRows(rows)

	// Act
	_, err := db.FindAll("")
	// Assert
	if err != nil {
		t.Errorf("test find all customers failed: %v", err)
	}
}

func TestCustomerRepositoryDb_FindAll_Failed(t *testing.T) {
	// Arrange
	client := setup()
	db := NewCustomerRepositoryDb(client)

	mock.ExpectQuery(
		"^SELECT \\* FROM `customers`").
		//WillReturnRows(rows)
		WillReturnError(errs.NewUnexpectedError("Unexpected database error"))
	// Act
	_, err := db.FindAll("")
	// Assert

	if want := errs.NewUnexpectedError("Unexpected database error"); !reflect.DeepEqual(err, want) {
		t.Errorf("validate save customer error failed:\ngot: %v\n want: %v\n", err, want)
	}
}
