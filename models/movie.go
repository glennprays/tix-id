package models

type Movie struct {
	ID          *int    `json:"id,omitempty"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    int     `json:"duration"`
	Rating      float32 `json:"rating"`
	ReleaseDate string  `json:"releaseDate"`
}

type MovieSchedules struct {
	Movie
	Schedules []Schedule
}

type MovieSchedulesResponse struct {
	Response
	MovieSchedules
}

type MovieSchedule struct {
	Movie
	Schedule Schedule
}

type MoviesResponse struct {
	Response
	Movies []Movie `json:"data"`
}
type MovieResponse struct {
	Response
	Movie Movie `json:"data"`
}
