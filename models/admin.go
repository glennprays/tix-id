package models

type Admin struct {
	ID       int    `json:"ID"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Phone    string `json:"Phone"`
	NIK      string `json:"NIK"`
}
