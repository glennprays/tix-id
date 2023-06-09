package models

type Seat struct {
	ID           int    `json:"id"`
	Row          string `json:"row,omitempty"`
	Number       string `json:"number,omitempty"`
	Availability *bool  `json:"availability,omitempty"`
}

type SeatRow struct {
	Row   string `json:"row"`
	Count int    `json:"count"`
}
