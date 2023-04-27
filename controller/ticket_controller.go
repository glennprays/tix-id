package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"tix-id/config"
	"tix-id/middleware"
	"tix-id/models"

	"github.com/gin-gonic/gin"
)

// GetPayments godoc
// @Summary Get Payments
// @Description Get all payments
// @Tags Customer
// @Accept json
// @Produce json
// @Param customerId path int true "Customer ID"
// @Success 200 {object} models.TicketsResponse
// @Router /customer/{customerId}/tickets [get]
func GetTickets(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()
	// TODO:get by customer id and verify with id in cookies
	customerIdParam, err := strconv.Atoi(c.Param("customerId"))
	customerId, _, _ := middleware.GetUserIdAndRoleFromCookie(c)
	if customerIdParam != int(customerId) {
		response := models.Response{
			Status:  200,
			Message: "The user id didn't matched",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// dummy data
	// customerId := customerIdParam

	query := "select count(*) from ticket where customer_id = ?"
	var count int
	err = db.QueryRow(query, customerId).Scan(&count)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if count < 1 {
		response := models.Response{
			Status:  404,
			Message: "You have no tickets!",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	// get data
	query = "SELECT tc.id, se.id, se.row, se.seat_number, p.id, p.amount, p.payment_status, s.id, s.price, s.show_time, m.id, m.title, m.description, m.duration, m.rating, m.release_date, b.id, b.name, b.address, t.id, t.name from ticket tc join seat se on se.id = tc.seat_id join payment p on p.id = tc.payment_id join schedule s on s.id = tc.schedule_id join movie m on m.id = s.movie_id join theatre t on t.id = s.theatre_id join branch b on b.id = t.branch_id where tc.customer_id = ?"
	rows, err := db.Query(query, customerId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ticket models.Ticket
	var tickets []models.Ticket
	var seat models.Seat
	var payment models.Payment
	var schedule models.ScheduleTicket
	var movie models.Movie
	var branch models.BranchTheatre
	var theatre models.Theatre
	for rows.Next() {
		if err := rows.Scan(&ticket.ID, &seat.ID, &seat.Row, &seat.Number, &payment.ID, &payment.Amount, &payment.Status, &schedule.ID, &schedule.Price, &schedule.Showtime, &movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.Rating, &movie.ReleaseDate, &branch.ID, &branch.Name, &branch.Address, &theatre.ID, &theatre.Name); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			branch.Theatre = theatre
			schedule.Branch = &branch
			schedule.Movie = &movie
			ticket.Schedule = schedule
			ticket.Payment = payment
			ticket.Seat = seat
			tickets = append(tickets, ticket)
		}
	}

	responseData := models.TicketsResponse{
		Response: models.Response{
			Status:  200,
			Message: "Tickets retrieved successfully",
		},
		Tickets: tickets,
	}

	c.JSON(http.StatusOK, responseData)
}

// CreateTicket godoc
// @Summary Create Ticket
// @Description Create Ticket by shedule
// @Tags Customer
// @Accept json
// @Produce json
// @Param customerId path int true "Customer ID"
// @Param body body models.ScheduleTicket true "Schedule Detail"
// @Success 200 {object} models.TicketResponse
// @Router /customer/{customerId}/tickets [post]
func CreateTicket(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()
	// TODO:get by customer id and verify with id in cookies
	customerIdParam, err := strconv.Atoi(c.Param("customerId"))
	// customerId, _, _ := middleware.GetUserIdAndRoleFromCookie(c)
	// if customerIdParam != int(customerId) {
	// 	response := models.Response{
	// 		Status:  200,
	// 		Message: "The user id didn't matched",
	// 	}
	// 	c.JSON(http.StatusOK, response)
	// 	return
	// }

	// dummy data
	customerId := customerIdParam
	var schedule models.ScheduleTicket
	if err := c.ShouldBindJSON(&schedule); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// checking schedule
	var movie models.Movie
	var theatre models.Theatre
	error := db.QueryRow("select movie_id, theatre_id, show_time from schedule where id = ?", schedule.ID).Scan(&movie.ID, &theatre.ID, &schedule.Showtime)
	if error != nil {
		log.Println(error)
		if error == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Schledule is not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	// get movie data
	error = db.QueryRow("select title, description, duration, rating, release_date from movie where id = ?", movie.ID).Scan(&movie.Title, &movie.Description, &movie.Duration, &movie.Rating, &movie.ReleaseDate)

	if error != nil {
		log.Println("error disini")
		log.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}
	schedule.Movie = &movie

	// get branch data
	var branch models.BranchTheatre
	branch.Theatre = theatre
	error = db.QueryRow("select b.name, b.address, t.name from branch b join theatre t on t.branch_id = b.id where t.id = ?", branch.Theatre.ID).Scan(&branch.Name, &branch.Address, &branch.Theatre.Name)
	if error != nil {
		log.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	var seat models.Seat
	seat.ID = schedule.Seat.ID

	// verify  seat is not taken yet
	var count int
	error = db.QueryRow("select count(*) from ticket where seat_id = ? and schedule_id = ?", seat.ID, schedule.ID).Scan(&count)
	if error != nil {
		log.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	if count > 0 {
		response := models.Response{
			Status:  200,
			Message: "the seat is taken",
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// get seat data and verify the seat is matched with the schdule
	error = db.QueryRow("select row, seat_number from seat where id = ? and schedule_id = ?", seat.ID, schedule.ID).Scan(&seat.Row, &seat.Number)
	if error != nil {
		if error == sql.ErrNoRows {
			response := models.Response{
				Status:  404,
				Message: "the seat is not match with the schedule",
			}
			c.JSON(http.StatusNotFound, response)
			return
		}
		log.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	// make a new payment
	res, errQuery := db.Exec("insert into payment(amount, payment_status) values ((select price from schedule where id = ? limit 1),'pending')", schedule.ID)
	if errQuery != nil {
		log.Println(errQuery)
		c.JSON(http.StatusBadRequest, gin.H{"error": errQuery.Error()})
		return
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var payment models.Payment
	payment.ID = int(lastInsertID)

	// get payment data
	error = db.QueryRow("select amount, payment_status from payment where id = ?", payment.ID).Scan(&payment.Amount, &payment.Status)
	if error != nil {
		log.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	// create new ticket
	res, errQuery = db.Exec("insert into ticket(customer_id, schedule_id, seat_id, payment_id) values (?,?,?,?)",
		customerId,
		schedule.ID,
		schedule.Seat.ID,
		payment.ID,
	)
	lastInsertID, err = res.LastInsertId()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule.Seat = nil
	var ticket models.Ticket
	ticket.ID = int(lastInsertID)
	ticket.Payment = payment
	ticket.Seat = seat
	ticket.Schedule = schedule

	responseData := models.TicketResponse{
		Response: models.Response{
			Status:  200,
			Message: "Ticket created successfully",
		},
		Ticket: ticket,
	}

	c.JSON(http.StatusOK, responseData)
}

// GetPayment godoc
// @Summary Get Payment
// @Description Get payment
// @Tags Customer
// @Accept json
// @Produce json
// @Param customerId path int true "Customer ID"
// @Param ticketId path int true "Customer ID"
// @Success 200 {object} models.TicketResponse
// @Router /customer/{customerId}/tickets/{ticketId} [get]
func GetTicket(c *gin.Context) {
	ticketId := c.Param("ticketId")

	db := config.ConnectDB()
	defer db.Close()
	// TODO:get by customer id and verify with id in cookies
	customerIdParam, err := strconv.Atoi(c.Param("customerId"))
	// customerId, _, _ := middleware.GetUserIdAndRoleFromCookie(c)
	// if customerIdParam != int(customerId) {
	// 	response := models.Response{
	// 		Status:  200,
	// 		Message: "The user id didn't matched",
	// 	}
	// 	c.JSON(http.StatusOK, response)
	// 	return
	// }

	// dummy data
	customerId := customerIdParam

	// get data
	query := "SELECT tc.id, se.id, se.row, se.seat_number, p.id, p.amount, p.payment_status, s.id, s.price, s.show_time, m.id, m.title, m.description, m.duration, m.rating, m.release_date, b.id, b.name, b.address, t.id, t.name from ticket tc join seat se on se.id = tc.seat_id join payment p on p.id = tc.payment_id join schedule s on s.id = tc.schedule_id join movie m on m.id = s.movie_id join theatre t on t.id = s.theatre_id join branch b on b.id = t.branch_id where tc.id = ? and tc.customer_id = ?"
	var ticket models.Ticket
	var seat models.Seat
	var payment models.Payment
	var schedule models.ScheduleTicket
	var movie models.Movie
	var branch models.BranchTheatre
	var theatre models.Theatre
	err = db.QueryRow(query, ticketId, customerId).Scan(&ticket.ID, &seat.ID, &seat.Row, &seat.Number, &payment.ID, &payment.Amount, &payment.Status, &schedule.ID, &schedule.Price, &schedule.Showtime, &movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.Rating, &movie.ReleaseDate, &branch.ID, &branch.Name, &branch.Address, &theatre.ID, &theatre.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Status:  404,
				Message: "the ticket is not found!",
			}
			c.JSON(http.StatusNotFound, response)
			return
		}
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	branch.Theatre = theatre
	schedule.Branch = &branch
	schedule.Movie = &movie
	ticket.Schedule = schedule
	ticket.Payment = payment
	ticket.Seat = seat

	responseData := models.TicketResponse{
		Response: models.Response{
			Status:  200,
			Message: "Ticket retrieved successfully",
		},
		Ticket: ticket,
	}

	c.JSON(http.StatusOK, responseData)
}

// ConfirmPayment godoc
// @Summary Confirm Payment
// @Description Confirm payment with payment_id.
// @Tags Customer
// @Param customerId path int true "Customer ID"
// @Param ticketId path string true "payment id"
// @Accept json
// @Produce json
// @Success 200 {object} models.TicketResponse
// @Router /customer/{customerId}/tickets/{ticketId}/payment [post]
func ConfirmPayment(c *gin.Context) {

	// customerId := c.Param("customerId")

	db := config.ConnectDB()
	defer db.Close()
	// TODO:get by customer id and verify with id in cookies
	customerIdParam, err := strconv.Atoi(c.Param("customerId"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ticketId, err := strconv.Atoi(c.Param("ticketId"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// customerId, _, _ := middleware.GetUserIdAndRoleFromCookie(c)
	// if customerIdParam != int(customerId) {
	// 	response := models.Response{
	// 		Status:  200,
	// 		Message: "The user id didn't matched",
	// 	}
	// 	c.JSON(http.StatusOK, response)
	// 	return
	// }

	// dummy
	customerId := customerIdParam

	var ticket models.Ticket
	ticket.ID = ticketId

	// check if the ticket id matched with customer
	var payment models.Payment
	var count int
	var paymentID sql.NullInt64
	error := db.QueryRow("select count(*), p.id from ticket t join payment p on p.id = t.payment_id where t.id = ? and t.customer_id = ? and p.payment_status = 'pending'", ticket.ID, customerId).Scan(&count, &paymentID)
	if error != nil {
		log.Println(error)
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	if count != 1 || !paymentID.Valid {
		response := models.Response{
			Status:  404,
			Message: "either ticket and customer didn't matched or ticket is paid",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	payment.ID = int(paymentID.Int64)

	// set payment into completed
	res, err := db.Exec("update payment set payment_status = 'completed' where id = ?", payment.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		response := models.Response{
			Status:  404,
			Message: "payment not found",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	// get seat data
	var seat models.Seat
	error = db.QueryRow("select s.id, s.row, s.seat_number from seat s join ticket t on s.id = t.schedule_id where t.id = ?", ticket.ID).Scan(&seat.ID, &seat.Row, &seat.Number)
	if error != nil {
		log.Println(error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	ticket.Seat = seat

	// get payment data
	error = db.QueryRow("select amount, payment_status from payment where id = ?", payment.ID).Scan(&payment.Amount, &payment.Status)
	if error != nil {
		log.Println(error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	ticket.Payment = payment

	// get schedule data
	var schedule models.ScheduleTicket
	var movie models.Movie
	var theatre models.Theatre
	var branch models.BranchTheatre
	error = db.QueryRow("select m.id, m.title, m.description, m.duration, m.rating, m.release_date, t.id, t.name, b.id, b.name, b.address, s.id, s.show_time, s.price from movie m join schedule s on s.movie_id = m.id join theatre t on t.id = s.theatre_id join branch b on b.id = t.branch_id join ticket tc on tc.schedule_id = s.id where tc.id = ?", ticket.ID).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.Rating, &movie.ReleaseDate, &theatre.ID, &theatre.Name, &branch.ID, &branch.Name, &branch.Address, &schedule.ID, &schedule.Showtime, &schedule.Price)
	if error != nil {
		log.Println(error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	branch.Theatre = theatre
	schedule.Branch = &branch
	schedule.Movie = &movie
	ticket.Schedule = schedule

	var customer models.Customer
	// Check if schedule exists in database
	if err := db.QueryRow("SELECT name, email FROM customer WHERE id = ?", customerIdParam).Scan(
		&customer.Name,
		&customer.Email,
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}
	log.Println(ticket.Seat.Number)
	schedule.Seat = &models.Seat{}
	schedule.Seat.Number = ticket.Seat.Number
	schedule.Seat.Number = ticket.Seat.Row
	content := GenerateEmail(customer, payment, schedule)
	SendEmail(content, customer.Email)

	responseData := models.TicketResponse{
		Response: models.Response{
			Status:  200,
			Message: "Payment Confirm successfully",
		},
		Ticket: ticket,
	}

	c.JSON(http.StatusOK, responseData)
}
