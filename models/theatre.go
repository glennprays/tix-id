package models

type Theatre struct {
	ID           int    `json:"ID"`
	Name         string `json:"Name"`
	Seats        []Seat `json:"Seats"`
	Availability bool   `json:"Availability"`
}
