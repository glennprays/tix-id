package models

type Schedule struct {
	ID           int    `json:"ID"`
	Price        int    `json:"Price"`
	Showtime     string `json:"Showtime"`
	Movie        Movie  `json:"Movie"`
	Seats        []Seat `json:"Seats"`
	Availability bool   `json:"Availability"`
}
