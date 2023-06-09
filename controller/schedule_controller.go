package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"tix-id/config"
	"tix-id/models"

	"github.com/gin-gonic/gin"
)

// GetSchedule godoc
// @Summary Get Schedule
// @Description Get a schedule by movie_id and schedule_id.
// @Tags Customer
// @Param movieId path string true "movie id"
// @Accept json
// @Produce json
// @Success 200 {object} models.MovieSchedulesResponse
// @Router /movies/{movieId}/schedules [get]
func GetSchedules(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	movieId, err := strconv.Atoi(c.Param("movieId"))
	var movie models.Movie
	// get movie data
	err = db.QueryRow("select id, title, description, duration, rating, release_date from movie where id = ?", movieId).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.Rating, &movie.ReleaseDate)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Status:  404,
				Message: "the movie is not found!",
			}
			c.JSON(http.StatusNotFound, response)
			return
		}
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check movie schedule exist
	var count int
	err = db.QueryRow("select count(*) from schedule where movie_id = ?", movie.ID).Scan(&count)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if count < 1 {
		response := models.Response{
			Status:  404,
			Message: "the movie have no any schedules!",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	// get Schedules
	var schedules []models.Schedule
	query := "select sc.id, sc.show_time, sc.price, t.id, t.name, b.id, b.name, b.address from schedule sc join theatre t on t.id = sc.theatre_id join branch b on b.id = t.branch_id where sc.movie_id = ?"
	rows, err := db.Query(query, movie.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var schedule models.Schedule
		var theatre models.Theatre
		var branch models.BranchTheatre
		if err := rows.Scan(&schedule.ID, &schedule.Showtime, &schedule.Price, &theatre.ID, &theatre.Name, &branch.ID, &branch.Name, &branch.Address); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		branch.Theatre = theatre
		schedule.Branch = branch
		schedules = append(schedules, schedule)
	}

	responseData := models.MovieSchedulesResponse{
		Response: models.Response{
			Status:  200,
			Message: "Schedules retrieved successfully",
		},
		MovieSchedules: models.MovieSchedules{
			// Movie:     movie,
			Schedules: schedules,
		},
	}
	c.JSON(http.StatusOK, responseData)
}

// GetSchedulesByMovieId godoc
// @Summary Get all schedules for a movie by ID
// @Description Get all schedules for a movie by ID
// @Tags Customer
// @Accept json
// @Produce json
// @Param movieId path string true "Movie ID"
// @Param scheduleId path string true "Schedule ID"
// @Success 200 {object} models.MovieScheduleResponse
// @Router /movies/{movieId}/schedules/{scheduleId} [get]
func GetSchedule(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	movieId := c.Param("movieId")
	scheduleId := c.Param("scheduleId")

	// get movie data
	var movie models.Movie
	// get movie data
	err := db.QueryRow("select id, title, description, duration, rating, release_date from movie where id = ?", movieId).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.Rating, &movie.ReleaseDate)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Status:  404,
				Message: "the movie is not found!",
			}
			c.JSON(http.StatusNotFound, response)
			return
		}
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get schedule data
	var schedule models.Schedule
	var theatre models.Theatre
	var branch models.BranchTheatre

	err = db.QueryRow("select sc.id, sc.show_time, sc.price, t.id, t.name, b.id, b.name, b.address from schedule sc join theatre t on t.id = sc.theatre_id join branch b on b.id = t.branch_id where sc.id = ?", scheduleId).Scan(&schedule.ID, &schedule.Showtime, &schedule.Price, &theatre.ID, &theatre.Name, &branch.ID, &branch.Name, &branch.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Status:  404,
				Message: "the schedule is not found!",
			}
			c.JSON(http.StatusNotFound, response)
			return
		}
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	branch.Theatre = theatre
	schedule.Branch = branch

	// verify schedule have any seats
	var count int
	err = db.QueryRow("select count(*) from seat where schedule_id = ?", schedule.ID).Scan(&count)

	if count < 1 {
		response := models.Response{
			Status:  404,
			Message: "the schedule have no any seats!",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	// get seats data
	var seats []models.Seat
	query := "select se.id, se.row, se.seat_number, IF(EXISTS (SELECT 1 FROM ticket t WHERE t.seat_id = se.id), 0, 1) AS availability from seat se where se.schedule_id = ?"
	rows, err := db.Query(query, schedule.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var availabilityInt int
		var seat models.Seat
		if err := rows.Scan(&seat.ID, &seat.Row, &seat.Number, &availabilityInt); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		availability := availabilityInt == 1
		seat.Availability = &availability
		seats = append(seats, seat)
	}
	schedule.Seats = &seats

	responseData := models.MovieScheduleResponse{
		Response: models.Response{
			Status:  200,
			Message: "Schedule retrieved successfully",
		},
		MovieSchedule: models.MovieSchedule{
			Schedule: schedule,
		},
	}
	c.JSON(http.StatusOK, responseData)
}

// CreateMovieSchedule godoc
// @Summary Create a movie schedule
// @Description Create a schedule for a movie
// @Tags Admin
// @Param movieId path string true "Movie ID"
// @Param body body models.Schedule true "Schedule details"
// @Success 201 {object} models.ScheduleResponse
// @Router /movies/{movieId}/schedules [post]
func CreateMovieSchedule(c *gin.Context) {
	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()

	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movieId, err := strconv.Atoi(c.Param("movieId"))
	log.Println("movieid: ", movieId)

	// Check if movie exists in database
	if schedule.Movie != nil {
		var movie models.Movie
		if err := db.QueryRow("SELECT id FROM movie WHERE id = ?", movieId).Scan(&movie.ID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
	}
	// Check if branch exists in database
	if schedule.Branch.ID != nil {
		var branch models.BranchTheatre
		if err := db.QueryRow("SELECT id FROM branch WHERE id = ?", schedule.Branch.ID).Scan(&branch.ID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Branch not found"})
			return
		}
	}
	log.Println("A: ", schedule.Price)
	log.Println(" B: ", schedule.Showtime)
	log.Println(" D: ", schedule.Branch.ID)
	// Insert schedule into database
	result, err := db.Exec("INSERT INTO schedule (price, show_time, movie_id, theatre_id) VALUES (?, ?, ?, ?)",
		schedule.Price, schedule.Showtime, movieId, schedule.Branch.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the ID of the inserted schedule
	scheduleID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	schedule.ID = int(scheduleID)

	responseData := models.ScheduleResponse{
		Response: models.Response{
			Status:  200,
			Message: "Schedule inserted successfully",
		},
		Schedule: schedule,
	}
	c.JSON(http.StatusOK, responseData)
}

// AddScheduleSeats godoc
// @Summary Adding a scheudules seats
// @Description Add seats for the schedule
// @Tags Admin
// @Param movieId path string true "Movie ID"
// @Param scheduleId path string true "Schedule ID"
// @Param body body []models.SeatRow true "List of Seat Row Detail"
// @Success 201 {object} models.Response
// @Router /movies/{movieId}/schedules/{scheduleId}/seats [post]
func AddScheduleSeats(c *gin.Context) {
	db := config.ConnectDB()
	defer db.Close()

	movieId := c.Param("movieId")
	scheduleId := c.Param("scheduleId")

	var seatRows []models.SeatRow
	if err := c.ShouldBindJSON(&seatRows); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if the schedule exist
	var count int
	err := db.QueryRow("select count(*) from schedule where id = ? and movie_id = ?", scheduleId, movieId).Scan(&count)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		response := models.Response{
			Status:  404,
			Message: "the schedule is not found!",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	for _, seat := range seatRows {
		// check if the row already exist
		var lastNumber int

		err = db.QueryRow("select seat_number from seat where row = ? and schedule_id = ? order by seat_number DESC LIMIT 1", seat.Row, scheduleId).Scan(&lastNumber)
		if err != nil {
			if err == sql.ErrNoRows {
				lastNumber = 1
			} else {

				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		if lastNumber != 1 {
			lastNumber++
		}

		for i := (0 + lastNumber); i <= seat.Count; i++ {
			_, err := db.Exec("insert into seat (row, seat_number, schedule_id) values (?, ?, ?)", seat.Row, i, scheduleId)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}

	responseData := models.Response{
		Status:  200,
		Message: "Seats inserted successfully!",
	}
	c.JSON(http.StatusOK, responseData)

}

// UpdateMovieSchedule godoc
// @Summary Update a movie schedule
// @Description Update a schedule for a movie
// @Tags Admin
// @Param movieId path string true "Movie ID"
// @Param scheduleId path string true "Schedule ID"
// @Param body body models.Schedule true "Schedule details"
// @Success 204 {object} models.ScheduleResponse
// @Router /movies/{movieId}/schedule/{scheduleId} [put]
func UpdateMovieSchedule(c *gin.Context) {
	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()

	movieId, err := strconv.Atoi(c.Param("movieId"))
	var schedule models.Schedule
	schedule.Branch = models.BranchTheatre{}
	schedule.Movie = &models.Movie{}
	schedule.Movie.ID = movieId
	schedule.Branch.Theatre = models.Theatre{}

	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("schedule.Branch.Theatre.ID: ", schedule.Branch.Theatre.ID)
	scheduleID, err := strconv.Atoi(c.Param("scheduleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}
	fmt.Println("theatre checkpoint 1: ", schedule.Branch.Theatre.ID)
	var count int
	err = db.QueryRow("select count(*) from ticket where schedule_id = ?", scheduleID).Scan(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule already has tickets, it cannot be changed"})
		return
	}
	var schedulee models.Schedule

	// Check if schedule exists in database
	if err := db.QueryRow("SELECT id, price, show_time, theatre_id FROM schedule WHERE id = ?", scheduleID).Scan(
		&schedulee.ID,
		&schedulee.Price,
		&schedulee.Showtime,
		&schedulee.Branch.ID,
	); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	// Update schedule in database
	log.Println("schedule.Price: ", schedule.Price)
	log.Println("schedule.Showtime: ", schedule.Showtime)
	log.Println("schedule.Movie.ID: ", schedule.Movie.ID)
	log.Println("scheduleID: ", scheduleID)
	log.Println("schedule.Branch.Theatre.ID: ", schedule.Branch.Theatre.ID)

	_, err = db.Exec("UPDATE schedule SET price = ?, show_time = ?, movie_id = ?, theatre_id = ? WHERE id = ?",
		schedule.Price, schedule.Showtime, schedule.Movie.ID, schedule.Branch.Theatre.ID, scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error1": err.Error()})
		return
	}
	schedule.Branch = models.BranchTheatre{}
	schedulee.Movie = &models.Movie{}
	schedulee.Branch = models.BranchTheatre{}
	if err := db.QueryRow("SELECT id, price, show_time, movie_id, theatre_id FROM schedule WHERE id = ?", scheduleID).Scan(
		&schedulee.ID,
		&schedulee.Price,
		&schedulee.Showtime,
		&schedulee.Movie.ID,
		&schedulee.Branch.ID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	responseData := models.ScheduleResponse{
		Response: models.Response{
			Status:  200,
			Message: "Schedule updated successfully",
		},
		Schedule: schedulee,
	}
	c.JSON(http.StatusOK, responseData)
}

// DeleteSchedule godoc
// @Summary Delete movie schedule
// @Description Delete a specific movie schedule by id
// @Tags Admin
// @Param movieId path string true "Movie ID"
// @Param scheduleId path string true "Schedule ID"
// @Produce json
// @Success 204 {object} models.Response
// @Router /movies/{movieId}/schedule/{scheduleId} [delete]
func DeleteSchedule(c *gin.Context) {
	scheduleID, err := strconv.Atoi(c.Param("scheduleId"))
	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()

	var count int
	var date time.Time
	err = db.QueryRow("select count(*) from ticket where schedule_id = ?", scheduleID).Scan(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule already has tickets, it cannot be changed"})
		err = db.QueryRow("select show_time from schedule where schedule_id = ?", scheduleID).Scan(&date)
		today := time.Now()
		older := today.After(date)
		if !older {
			return
		}
	}

	log.Println("scheduleId: ", scheduleID)

	// Delete the movie
	result, err := db.Exec("DELETE FROM schedule WHERE id=?", scheduleID)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	responseData := models.Response{
		Status:  http.StatusOK,
		Message: "Schedule deleted successfully",
	}
	c.JSON(http.StatusOK, responseData)
}
