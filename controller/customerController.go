package controller

import (
	"fmt"
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
// @Success 200 {object} models.Customer
// @Failure 400 {object} map[string]interface{}
// @Router /customer/registration [post]
func AddCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("berhasil log")

	c.JSON(http.StatusCreated, gin.H{"data": customer})
}

// LoginCustomer godoc
// @Summary Login Customer
// @Description Login Customer Account
// @Tags Customer
// @Accept json
// @Produce json
// @Success 200 {object} models.Customer
// @Failure 400 {object} map[string]interface{}
// @Router /customer/auth/login [post]
func LoginCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("berhasil log")

	c.JSON(http.StatusCreated, gin.H{"data": customer})
}

// GetCustomer godoc
// @Summary Get Customer
// @Description get Account
// @Tags Customer
// @Accept json
// @Produce json
// @Success 200 {object} models.Customer
// @Failure 400 {object} map[string]interface{}
// @Router /customer/{customer_id} [get]
func GetCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("berhasil log")

	c.JSON(http.StatusCreated, gin.H{"data": customer})
}

// PutCustomer godoc
// @Summary Put Customer
// @Description Put Customer Account
// @Tags Customer
// @Accept json
// @Produce json
// @Success 200 {object} models.Customer
// @Failure 400 {object} map[string]interface{}
// @Router /customer/{customer_id} [put]
func UpdateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("berhasil log")

	c.JSON(http.StatusCreated, gin.H{"data": customer})
}

// GetMovies godoc
// @Summary Get Movies
// @Description Get All Movies sorted by realse date
// @Tags Customer
// @Param show_time query string false "Filter movies by show time"
// @Param branch query string false "Filter movies by branch"
// @Param rating query string false "Filter movies by rating"
// @Accept json
// @Produce json
// @Success 200 {object} models.Movie
// @Failure 400 {object} map[string]interface{}
// @Router /movies [get]
func GetMovies(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("berhasil log")

	c.JSON(http.StatusCreated, gin.H{"data": movie})
}

// SearchMovies godoc
// @Summary Search movies by title and genre
// @Description Search movies by title and genre
// @Tags Customer
// @Accept json
// @Produce json
// @Param title query string false "Movie title to search"
// @Param genre query string false "Movie genre to search"
// @Success 200 {array} models.Movie
// @Failure 400 {object} map[string]string
// @Router /movies/search [get]
func SearchMovies(c *gin.Context) {
	// TODO: Implement logic to search movies by title and genre
	var movies []models.Movie

	// Return movies as response
	c.JSON(http.StatusOK, movies)
}

// GetMovieById godoc
// @Summary Get a movie by ID
// @Description Get a movie by ID
// @Tags Customer
// @Accept json
// @Produce json
// @Param movie_id path string true "Movie ID"
// @Success 200 {object} models.Movie
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /movies/{movie_id} [get]
func GetMovieById(c *gin.Context) {
	// TODO: Implement logic to fetch movie by ID
	var movie models.Movie

	// Return movie as response
	c.JSON(http.StatusOK, movie)
}

// GetSchedulesByMovieId godoc
// @Summary Get all schedules for a movie by ID
// @Description Get all schedules for a movie by ID
// @Tags Customer
// @Accept json
// @Produce json
// @Param movie_id path string true "Movie ID"
// @Success 200 {array} models.Schedule
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /movies/{movie_id}/schedules [get]
func GetSchedulesByMovieId(c *gin.Context) {
	// TODO: Implement logic to fetch schedules by movie ID
	var schedules []models.Schedule

	// Return schedules as response
	c.JSON(http.StatusOK, schedules)
}

// GetSchedule godoc
// @Summary Get Schedule
// @Description Get a schedule by movie_id and schedule_id.
// @Tags Customer
// @Param movie_id path string true "movie id"
// @Param schedule_id path string true "schedule id"
// @Accept json
// @Produce json
// @Success 200 {object} models.Schedule
// @Failure 400 {object} map[string]interface{}
// @Router /movie/{movie_id}/schedule/{schedule_id} [get]
func GetSchedule(c *gin.Context) {

	// Your logic here to fetch the schedule with given movieID and scheduleID
	// and return the result in JSON format.
	// Example:
	var schedule models.Schedule
	c.JSON(http.StatusOK, schedule)
}

// UpdateSchedule godoc
// @Summary Update Schedule
// @Description Update a schedule by movie_id and schedule_id.
// @Tags Customer
// @Param movie_id path string true "movie id"
// @Param schedule_id path string true "schedule id"
// @Accept json
// @Produce json
// @Success 200 {object} models.Schedule
// @Failure 400 {object} map[string]interface{}
// @Router /movie/{movie_id}/schedule/{schedule_id} [put]
func UpdateSchedule(c *gin.Context) {

	// Your logic here to update the schedule with given movieID and scheduleID
	// using the request body, and return the updated schedule in JSON format.
	// Example:
	var schedule models.Schedule
	c.JSON(http.StatusOK, schedule)
}

// ConfirmPayment godoc
// @Summary Confirm Payment
// @Description Confirm payment with payment_id.
// @Tags Customer
// @Param payment_id path string true "payment id"
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /customer/ticket/payment/{payment_id} [post]
func ConfirmPayment(c *gin.Context) {

	// Your logic here to confirm the payment with given paymentID
	// using the request body, and return the confirmation result in JSON format.
	// Example:
	result := map[string]string{
		"message": "Payment confirmed",
	}
	c.JSON(http.StatusOK, result)
}

// GetPayments godoc
// @Summary Get Payments
// @Description Get all payments
// @Tags Customer
// @Accept json
// @Produce json
// @Success 200 {object} []models.Payment
// @Failure 400 {object} map[string]interface{}
// @Router /customer/ticket/payments [get]
func GetPayments(c *gin.Context) {
	var payment models.Payment

	c.JSON(http.StatusOK, payment)
}
