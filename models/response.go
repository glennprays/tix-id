package models

type ResponseWithPayload struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Paging  *Paging     `json:"paging,omitempty"`
}

type BasicResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Paging struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
