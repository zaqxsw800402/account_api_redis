package dto

import (
	"net/http"
	"testing"
)

func Test_Transaction_Validate(t *testing.T) {
	request := TransactionRequest{
		TransactionType: "invalid",
	}

	appError := request.Validate()
	if appError.Message != "transaction type can only be deposit or withdrawal" {
		t.Error("Invalid transaction type while testing")
	}

	if appError.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid code while testing")
	}

}
