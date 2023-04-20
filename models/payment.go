package models

type Payment struct {
	ID     int           `json:"ID"`
	Amount int           `json:"Amount"`
	Status PaymentStatus `json:"Status"`
}

type PaymentStatus string

const (
	Pending   PaymentStatus = "Pending"
	Completed PaymentStatus = "Completed"
	Expired   PaymentStatus = "Expired"
)
