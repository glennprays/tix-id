package models

type Payment struct {
	ID     int    `json:"ID"`
	Amount int    `json:"Amount"`
	Status string `json:"Status"`
}
