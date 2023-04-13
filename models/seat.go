package models

type Seat struct {
	ID           int    `json:"ID"`
	Row          string `json:"Row"`
	Number       string `json:"Number"`
	Availability bool   `json:"Availability"`
}
