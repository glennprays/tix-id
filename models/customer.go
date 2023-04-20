package models

type Customer struct {
	ID       int     `json:"ID"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password *string `json:"password,omitempty"`
	Phone    string  `json:"phone"`
}

type CustomerResponse struct {
	Response
	Customer Customer `json:"data"`
}
