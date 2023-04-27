package models

type Theatre struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

type TheatresResponse struct {
	Response
	Theatres []Theatre `json:"data"`
}

type TheatreResponse struct {
	Response
	Theatre Theatre `json:"data"`
}
