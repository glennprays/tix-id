package models

import "time"

type Schedule struct {
	ID       int       `json:"id"`
	Price    int       `json:"price"`
	Showtime time.Time `json:"showtime"`
	Movie    Movie     `json:"movie"`
	Branch   Branch    `json:"branch"`
	Seats    []Seat    `json:"seats"`
}
