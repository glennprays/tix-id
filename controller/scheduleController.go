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
	// movieId := c.Query("movieId")
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseData := models.ScheduleResponse{
		Response: models.Response{
			Status:  200,
			Message: "Schedule inserted successfully",
		},
		Schedule: schedule,
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
	// movieId := c.Query("movieId")
	// scheduleId := c.Query("scheduleId")
	var schedule models.Schedule
	responseData := models.ScheduleResponse{
		Response: models.Response{
			Status:  200,
			Message: "Schedule updated successfully",
		},
		Schedule: schedule,
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
	// movieId := c.Query("movieId")
	// scheduleId := c.Query("scheduleId")
	responseData := models.Response{
		Status:  200,
		Message: "Schedule deleted successfully",
	}
	c.JSON(http.StatusOK, responseData)
}
