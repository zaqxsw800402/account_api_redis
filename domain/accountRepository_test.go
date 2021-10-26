package domain

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
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

func TestSave_Success(t *testing.T) {
	client := setup()
	//defer client.
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
	fmt.Println("err: ", err)
}
