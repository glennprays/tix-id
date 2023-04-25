package models

type Ticket struct {
	ID       int            `json:"id"`
	Seat     Seat           `json:"seat"`
	Schedule ScheduleTicket `json:"schedule"`
	Payment  Payment        `json:"payment"`
}

type TicketResponse struct {
	Response
	Ticket Ticket `json:"data"`
}
type TicketsResponse struct {
	Response
	Tickets []Ticket `json:"data"`
}
