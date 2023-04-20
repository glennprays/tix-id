package models

type Ticket struct {
	ID       int      `json:"id"`
	Seat     Seat     `json:"name"`
	Schedule Schedule `json:"address"`
	Payment  Payment  `json:"theatres"`
}
