package controller

import (
	"net/http"
	"tix-id/models"

	"github.com/gin-gonic/gin"
)

// GetPayments godoc
// @Summary Get Payments
// @Description Get all payments
// @Tags Customer/Payment
// @Accept json
// @Produce json
// @Param customerId path int true "Customer ID"
// @Success 200 {object} models.PaymentsResponse
// @Router /customer/{customerId}/ticket/payments [get]
func GetPayments(c *gin.Context) {
	// customerId := c.Param("customerId")

	var payments []models.Payment
	responseData := models.PaymentsResponse{
		Response: models.Response{
			Status:  200,
			Message: "Payments retrieved successfully",
		},
		Payments: payments,
	}

	c.JSON(http.StatusOK, responseData)
}

// GetPayment godoc
// @Summary Get Payment
// @Description Get payment
// @Tags Customer/Payment
// @Accept json
// @Produce json
// @Param customerId path int true "Customer ID"
// @Param paymentId path int true "Customer ID"
// @Success 200 {object} models.PaymentResponse
// @Router /customer/{customerId}/ticket/payments/{paymentId} [get]
func GetPayment(c *gin.Context) {
	// customerId := c.Param("customerId")
	// paymentId := c.Param("paymentId")

	var payment models.Payment
	responseData := models.PaymentResponse{
		Response: models.Response{
			Status:  200,
			Message: "Payment retrieved successfully",
		},
		Payment: payment,
	}

	c.JSON(http.StatusOK, responseData)
}

// ConfirmPayment godoc
// @Summary Confirm Payment
// @Description Confirm payment with payment_id.
// @Tags Customer/Payment
// @Param payment_id path string true "payment id"
// @Accept json
// @Produce json
// @Success 200 {object} models.PaymentResponse
// @Router /customer/{customerId}/ticket/payments/{payment_id} [post]
func ConfirmPayment(c *gin.Context) {

	// customerId := c.Param("customerId")
	// paymentId := c.Param("paymentId")

	var payment models.Payment
	responseData := models.PaymentResponse{
		Response: models.Response{
			Status:  200,
			Message: "Payment retrieved successfully",
		},
		Payment: payment,
	}

	c.JSON(http.StatusOK, responseData)
}
