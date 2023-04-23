package models

type Payment struct {
	ID     int           `json:"id"`
	Amount int           `json:"amount"`
	Status PaymentStatus `json:"status"`
}

type PaymentStatus string

const (
	Pending   PaymentStatus = "Pending"
	Completed PaymentStatus = "Completed"
	Expired   PaymentStatus = "Expired"
)

type PaymentResponse struct {
	Response
	Payment Payment `json:"data"`
}
type PaymentsResponse struct {
	Response
	Payments []Payment `json:"data"`
}
