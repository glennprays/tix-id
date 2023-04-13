package models

type Ticket struct {
	ID       int      `json:"ID"`
	Seat     Seat     `json:"Name"`
	Schedule Schedule `json:"Address"`
	Payment  Payment  `json:"Theatres"`
}
