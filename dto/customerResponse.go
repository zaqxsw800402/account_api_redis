package dto

type CustomerResponse struct {
	Id          uint        `json:"customer_id"`
	Name        string      `json:"name"`
	City        string      `json:"city"`
	Zipcode     string      `json:"zipcode"`
	DateOfBirth string      `json:"date_of_birth"`
	Status      string      `json:"status"`
	Accounts    interface{} `json:"accounts"`
}

type CustomerRequest struct {
	Name        string `json:"name"`
	City        string `json:"city"`
	Zipcode     string `json:"zipcode"`
	DateOfBirth string `json:"date_of_birth"`
}
