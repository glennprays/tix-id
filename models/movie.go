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
	Schedules []Schedule `json:"data"`
}

type MovieSchedulesResponse struct {
	Response
	MovieSchedules
}

type MovieSchedule struct {
	Schedule Schedule `json:"data"`
}

type MovieScheduleResponse struct {
	Response
	MovieSchedule
}

type MoviesResponse struct {
	Response
	Movies []Movie `json:"data"`
}
type MovieResponse struct {
	Response
	Movie Movie `json:"data"`
}
