package dto

type UserRequest struct {
	//ID        int       `json:"id" `
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name" `
	Email     string `json:"email" `
	Password  string `json:"password" `
	//CreatedAt time.Time `json:"-" `
	//UpdatedAt time.Time `json:"-" `
}

type UserResponse struct {
	ID        int    `json:"id" `
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password" `
	//CreatedAt time.Time `json:"-" `
	//UpdatedAt time.Time `json:"-" `
}
