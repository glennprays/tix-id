package models

type Movie struct {
	ID          *int    `json:"id,omitempty"`
	Title       int     `json:"title"`
	Description string  `json:"description"`
	Duration    int     `json:"duration"`
	Rating      float32 `json:"rating"`
	ReleaseDate string  `json:"releaseDate"`
}

type MoviesResponse struct {
	Response
	Movies []Movie `json:"data"`
}
type MovieResponse struct {
	Response
	Movie Movie `json:"data"`
}
