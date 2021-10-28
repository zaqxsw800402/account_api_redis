package domain

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"red/errs"
	"reflect"
	"regexp"
	"testing"
)

var mock sqlmock.Sqlmock

func setup() *gorm.DB {
	//sqlmock
	var err error
	var db *sql.DB
	db, mock, err = sqlmock.New()
	if nil != err {
		log.Fatalf("Init sqlmock failed, err %v", err)
	}
	//gorm、sqlmock
	client, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	if nil != err {
		log.Fatalf("Init DB with sqlmock failed, err %v", err)
	}
	return client
}

func TestAccountRepositoryDb_Save_Success(t *testing.T) {
	client := setup()
	db := NewAccountRepositoryDb(client)
	//rows := sqlmock.NewRows([]string{"customer_id","opening_date","account_type","amount","status","account_id"}).
	//	AddRow(2,"2012-10-18","saving",6000,1,1)
	//prep := mock.ExpectPrepare("^INSERT INTO accounts*")
	//prep.ExpectExec().WithArgs(2,"2012-10-18 10:10:10","saving",6000,"1").
	//	WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta("INSERT INTO `accounts` (`customer_id`,`opening_date`,`account_type`,`amount`,`status`,`account_id`) VALUES (?,?,?,?,?,?)")).
		//`INSERT INTO "accounts" ("customer_id","opening_date","account_type","amount","status","account_id") VALUES (?,?,?,?,?,?)`).
		WithArgs(2, "2012-10-18 10:10:10", "saving", 6000.00, "1", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	//mock.ExpectQuery("^INSERT INTO accounts").WithArgs(1,2,"2012-10-18","saving",6000,"1").WillReturnRows(rows)

	account := Account{
		AccountId:   1,
		CustomerId:  2,
		OpeningDate: "2012-10-18 10:10:10",
		AccountType: "saving",
		Amount:      6000,
		Status:      "1",
	}
	_, err := db.Save(account)
	if err != nil {
		t.Error("test failed while save account into database")
	}
}

func TestAccountRepositoryDb_Save_Failed(t *testing.T) {
	// Arrange
	client := setup()
	db := NewAccountRepositoryDb(client)

	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta("INSERT INTO `accounts` (`customer_id`,`opening_date`,`account_type`,`amount`,`status`,`account_id`) VALUES (?,?,?,?,?,?)")).
		WithArgs(2, "2012-10-18 10:10:10", "saving", 6000, "1", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	account := Account{
		AccountId:   1,
		CustomerId:  2,
		OpeningDate: "2012-10-18 10:10:10",
		AccountType: "saving",
		Amount:      6000,
		Status:      "1",
	}
	_, err := db.Save(account)
	// Assert
	if want := errs.NewUnexpectedError("Unexpected error from database"); err == nil {
		t.Errorf("test failed while save account into database got: %v\nwant: %v", err, want)
	}
}

func TestAccountRepositoryDb_FindBy_Success(t *testing.T) {
	// Arrange
	client := setup()
	db := NewAccountRepositoryDb(client)
	rows := sqlmock.NewRows([]string{"customer_id", "opening_date", "account_type", "amount", "status", "account_id"}).
		AddRow(2, "2012-10-18", "saving", 6000, 1, 1)

	transactionRows := sqlmock.NewRows([]string{"transaction_id", "account_id", "amount", "transaction_type", "transaction_date"}).
		AddRow(1, 1, 6000, "deposit", "2012-10-18")

	mock.ExpectQuery(
		"^SELECT \\* FROM `Accounts` WHERE account_id = \\?").
		WillReturnRows(rows)

	mock.ExpectQuery(
		"^SELECT \\* FROM `transactions` WHERE `transactions`.`account_id` = \\?").
		WithArgs(1).
		WillReturnRows(transactionRows)

	// Act
	_, err := db.FindBy(1)
	// Assert
	if err != nil {
		t.Error("test failed while find account through id from database")
	}
	fmt.Println("error:", err)
}

func TestAccountRepositoryDb_FindBy_Failed(t *testing.T) {
	// Arrange
	client := setup()
	db := NewAccountRepositoryDb(client)

	//rows := sqlmock.NewRows([]string{"customer_id","opening_date","account_type","amount","status","account_id"}).
	//	AddRow(2,"2012-10-18","saving",6000,1,1)

	transactionRows := sqlmock.NewRows([]string{"transaction_id", "account_id", "amount", "transaction_type", "transaction_date"}).
		AddRow(1, 1, 6000, "deposit", "2012-10-18")

	mock.ExpectQuery(
		"^SELECT \\* FROM `Accounts` WHERE account_id = \\?").
		//WillReturnRows(rows)
		WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery(
		"^SELECT \\* FROM `transactions` WHERE `transactions`.`account_id` = \\?").
		WithArgs(1).
		WillReturnRows(transactionRows)

	// Act
	_, err := db.FindBy(1)
	// Assert
	if want := errs.NewNotFoundError("Account not found"); !reflect.DeepEqual(err, want) {
		t.Errorf("test failed while get account from database got: %v\nwant: %v", err, want)
	}
	fmt.Println("error ", err)
}

func TestAccountRepositoryDb_SaveTransaction_Success(t *testing.T) {
	// Arrange
	client := setup()
	db := NewAccountRepositoryDb(client)

	mock.ExpectBegin()

	mock.ExpectExec(
		regexp.QuoteMeta("INSERT INTO `transactions` (`account_id`,`amount`,`transaction_type`,`transaction_date`,`transaction_id`) VALUES (?,?,?,?,?)")).
		WithArgs(1, 6000.00, "deposit", "2012-10-18", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(
		regexp.QuoteMeta("UPDATE `accounts` SET `amount`=amount + ? WHERE `account_id` = ?")).
		WithArgs(6000.00, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	rows := sqlmock.NewRows([]string{"customer_id", "opening_date", "account_type", "amount", "status", "account_id"}).
		AddRow(1, "2012-10-18", "saving", 12000, 1, 1)

	mock.ExpectQuery(
		"^SELECT \\* FROM `Accounts` WHERE account_id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	transactionRows := sqlmock.NewRows([]string{"transaction_id", "account_id", "amount", "transaction_type", "transaction_date"}).
		AddRow(1, 1, 6000, "deposit", "2012-10-18")

	mock.ExpectQuery(
		"^SELECT \\* FROM `transactions` WHERE `transactions`.`account_id` = \\?").
		WithArgs(1).
		WillReturnRows(transactionRows)

	transaction := Transaction{
		TransactionId:   1,
		AccountId:       1,
		Amount:          6000,
		TransactionType: "deposit",
		TransactionDate: "2012-10-18",
	}

	// Act
	_, err := db.SaveTransaction(transaction)

	// Assert
	if err != nil {
		t.Error("test failed while save transaction")
	}
}
