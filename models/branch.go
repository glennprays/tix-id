package models

type Branch struct {
	ID       *int       `json:"id,omitempty"`
	Name     int        `json:"name"`
	Address  string     `json:"address"`
	Theatres *[]Theatre `json:"theatres,omitempty"`
}

type BranchTheatre struct {
	ID      *int    `json:"id,omitempty"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Theatre Theatre `json:"theatre"`
}

type BranchesResponse struct {
	Response
	Branches []Branch `json:"data"`
}
type BranchResponse struct {
	Response
	Branch Branch `json:"data"`
}
