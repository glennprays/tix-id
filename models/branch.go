package models

type Branch struct {
	ID       int       `json:"ID"`
	Name     int       `json:"Name"`
	Address  string    `json:"Address"`
	Theatres []Theatre `json:"Theatres"`
}
