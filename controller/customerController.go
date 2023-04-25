package controller

import (
	"net/http"
	"tix-id/models"

	"github.com/gin-gonic/gin"
)

// CreateCustomer godoc
// @Summary Create New Customer
// @Description Create New Customer Account
// @Tags Customer
// @Accept json
// @Produce json
// @Param body body models.Customer true "Customer details"
// @Success 200 {object} models.CustomerResponse
// @Router /customer/registration [post]
func AddCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseData := models.CustomerResponse{
		Response: models.Response{
			Status:  200,
			Message: "Registration successful",
		},
		Customer: customer,
	}

	c.JSON(http.StatusCreated, responseData)
}

// LoginCustomer godoc
// @Summary Login Customer
// @Description Login Customer Account
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Login details"
// @Success 200 {object} models.CustomerResponse
// @Router /customer/auth/login [post]
func LoginCustomer(c *gin.Context) {
	var login models.LoginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer models.Customer
	responseData := models.CustomerResponse{
		Response: models.Response{
			Status:  200,
			Message: "Login successful",
		},
		Customer: customer,
	}

	c.JSON(http.StatusCreated, responseData)
}

// GetCustomer godoc
// @Summary Get Customer
// @Description get Account
// @Tags Customer
// @Accept json
// @Produce json
// @Param customerId path int true "Customer ID"
// @Success 200 {object} models.CustomerResponse
// @Router /customer/{customerId}/profile [get]
func GetCustomer(c *gin.Context) {

	// customerId := c.Param("customerId")

	var customer models.Customer
	responseData := models.CustomerResponse{
		Response: models.Response{
			Status:  200,
			Message: "Customer retrieved successfully",
		},
		Customer: customer,
	}

	c.JSON(http.StatusCreated, responseData)
}

// PutCustomer godoc
// @Summary Update Customer
// @Description Put Customer Account
// @Tags Customer
// @Accept json
// @Produce json
// @Param customerId path int true "Customer ID"
// @Param body body models.Customer true "Customer details"
// @Success 200 {object} models.CustomerResponse
// @Router /customer/{customerId}/profile [put]
func UpdateCustomer(c *gin.Context) {
	// customerId := c.Param("customerId")

	var customer models.Customer
	responseData := models.CustomerResponse{
		Response: models.Response{
			Status:  200,
			Message: "Customer updated successfully",
		},
		Customer: customer,
	}

	c.JSON(http.StatusCreated, responseData)
}
