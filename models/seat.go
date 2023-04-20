package models

type Seat struct {
	ID           int    `json:"id"`
	Row          string `json:"row"`
	Number       string `json:"number"`
	Availability bool   `json:"availability"`
}
