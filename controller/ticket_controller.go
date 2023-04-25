package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"tix-id/config"
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
func GetPayments(c *gin.Context) {
	// customerId := c.Param("customerId")

	var tickets []models.Ticket
	responseData := models.TicketsResponse{
		Response: models.Response{
			Status:  200,
			Message: "Payments retrieved successfully",
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
func GetPayment(c *gin.Context) {
	// customerId := c.Param("customerId")
	// ticketId := c.Param("ticketId")

	var ticket models.Ticket
	responseData := models.TicketResponse{
		Response: models.Response{
			Status:  200,
			Message: "Payment retrieved successfully",
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
	// ticketId := c.Param("ticketId")

	var ticket models.Ticket
	responseData := models.TicketResponse{
		Response: models.Response{
			Status:  200,
			Message: "Payment Confirm successfully",
		},
		Ticket: ticket,
	}

	c.JSON(http.StatusOK, responseData)
}
