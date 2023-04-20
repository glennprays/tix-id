package models

type Customer struct {
	ID       int     `json:"ID"`
	Name     string  `json:"Name"`
	Username string  `json:"Username"`
	Email    string  `json:"Email"`
	Password *string `json:"Password,omitempty"`
	Phone    string  `json:"Phone"`
}
