package dto

type CustomerResponse struct {
	Id          uint   `json:"customer_id"`
	Name        string `json:"name"`
	City        string `json:"city"`
	Zipcode     string `json:"zipcode"`
	DateOfBirth string `json:"date_of_birth"`
	Status      string `json:"status"`
	//Accounts    interface{} `json:"accounts"`
}

type CustomerRequest struct {
	//Id          int    `json:"customer_id"`
	Name string `json:"name" binding:"required"`
	City string `json:"city" binding:"required"`
	//Zipcode     string `json:"zipcode" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
}
