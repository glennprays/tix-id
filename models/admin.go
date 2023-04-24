package models

type Admin struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password *string `json:"password,omitempty"`
	Phone    string  `json:"phone"`
	NIK      string  `json:"NIK"`
}

type AdminResponse struct {
	Response
	Admin Admin `json:"data"`
}
