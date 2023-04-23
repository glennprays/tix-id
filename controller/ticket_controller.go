package controller

import (
	"net/http"
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

// GetPayment godoc
// @Summary Get Payment
// @Description Get payment
// @Tags Customer
// @Accept json
// @Produce json
// @Param customerId path int true "Customer ID"
// @Param ticketId path int true "Customer ID"
// @Success 200 {object} models.TicketResponse
// @Router /customer/{customerId}/ticket/{ticketId} [get]
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
// @Router /customer/{customerId}/ticket/{ticketId}/payment [post]
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
