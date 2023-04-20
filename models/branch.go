package models

type Branch struct {
	ID       *int      `json:"id,omitempty"`
	Name     int       `json:"name"`
	Address  string    `json:"address"`
	Theatres []Theatre `json:"theatres"`
}

type BranchesResponse struct {
	Response
	Branches []Branch `json:"data"`
}
type BranchResponse struct {
	Response
	Branch Branch `json:"data"`
}
