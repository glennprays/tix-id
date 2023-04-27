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
	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := db.Exec("INSERT INTO customer (username, password,name,email,phone) VALUES (?, ?,?,?,?)", customer.Username, customer.Password, customer.Name, customer.Email, customer.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Get the ID of the inserted customer
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ID of inserted branch"})
		return
	}

	// Set the ID of the customer to the inserted ID
	customer.ID = int(id)

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
	db := config.ConnectDB()
	defer db.Close()

	var login models.LoginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := db.QueryRow("select id, username, password, name, email, phone from customer where email = ? and password = ?",
		login.Email,
		login.Password)

	var customer models.Customer
	if err := row.Scan(&customer.ID, &customer.Username, &customer.Password, &customer.Name, &customer.Email, &customer.Phone); err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	} else {
		middleware.CreateToken(c, uint(customer.ID), "customer", 3600)

		responseData := models.CustomerResponse{
			Response: models.Response{
				Status:  200,
				Message: "Login successful",
			},
			Customer: customer,
		}

		c.JSON(http.StatusCreated, responseData)
	}
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
	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()
	// customerId := c.Param("customerId")
	customerId, err := strconv.Atoi(c.Param("customerId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer from database"})
		return
	}
	var customer models.Customer
	err = db.QueryRow("Select username,name,email,phone from customer where id =?", customerId).Scan(&customer.Username, &customer.Name, &customer.Email, &customer.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Status:  404,
				Message: "the branch is not found!",
			}
			c.JSON(http.StatusNotFound, response)
			return
		}
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer.ID = customerId

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
	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()
	// customerId := c.Param("customerId")
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customerId, err := strconv.Atoi(c.Param("customerId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer from database"})
		return
	}

	result, err := db.Exec("UPDATE customer SET username=?,password=?,name=?,email=?,phone=? WHERE id=?", customer.Username, customer.Password, customer.Name, customer.Email, customer.Phone, customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	customer.ID = customerId

	responseData := models.CustomerResponse{
		Response: models.Response{
			Status:  200,
			Message: "Customer updated successfully",
		},
		Customer: customer,
	}

	c.JSON(http.StatusCreated, responseData)
}
