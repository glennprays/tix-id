package models

type Movie struct {
	ID          int     `json:"ID"`
	Title       int     `json:"Title"`
	Description string  `json:"Description"`
	Duration    int     `json:"Duration"`
	Rating      float32 `json:"Rating"`
	ReleaseDate string  `json:"ReleaseDate"`
}
