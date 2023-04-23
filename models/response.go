package models

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Paging struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
