package models

type Payment struct {
	ID     int           `json:"id"`
	Amount float64       `json:"amount"`
	Status PaymentStatus `json:"status"`
}

type PaymentStatus string

const (
	Pending   PaymentStatus = "pending"
	Completed PaymentStatus = "completed"
	Failed    PaymentStatus = "failed"
)

type PaymentResponse struct {
	Response
	Payment Payment `json:"data"`
}
type PaymentsResponse struct {
	Response
	Payments []Payment `json:"data"`
}
