package dto

type CustomerResponse struct {
	Id          uint   `json:"customer_id"`
	Name        string `json:"name"`
	City        string `json:"city"`
	DateOfBirth string `json:"date_of_birth"`
	Status      string `json:"status"`
}

type CustomerRequest struct {
	//Id          int    `json:"customer_id"`
	Name        string `json:"name" binding:"required"`
	City        string `json:"city" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
}
