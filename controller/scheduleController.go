package controller

import (
	"net/http"
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
// @Success 200 {object} models.SchedulesResponse
// @Router /movies/{movieId}/schedules [get]
func GetSchedules(c *gin.Context) {
	// movieId := c.Query("movieId")
	var schedules []models.Schedule

	responseData := models.SchedulesResponse{
		Response: models.Response{
			Status:  200,
			Message: "Schedules retrieved successfully",
		},
		Schedules: schedules,
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
// @Success 200 {object} models.ScheduleResponse
// @Router /movies/{movieId}/schedules/{scheduleId} [get]
func GetSchedule(c *gin.Context) {
	// movieId := c.Query("movieId")
	// scheduleId := c.Query("scheduleId")
	var schedule models.Schedule

	responseData := models.ScheduleResponse{
		Response: models.Response{
			Status:  200,
			Message: "Schedule retrieved successfully",
		},
		Schedule: schedule,
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
