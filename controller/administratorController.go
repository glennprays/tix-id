package controller

import (
	"net/http"
	"tix-id/models"

	"github.com/gin-gonic/gin"
)

// CreateMovie godoc
// @Summary Create Movie
// @Description Create a new movie
// @Tags Admin
// @Accept json
// @Produce json
// @Param movie body models.Movie true "Movie object"
// @Success 201 {object} models.Movie
// @Failure 400 {object} map[string]interface{}
// @Router /movie [post]
func CreateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

// UpdateMovie godoc
// @Summary Update Movie
// @Description Update an existing movie
// @Tags Admin
// @Accept json
// @Produce json
// @Param movie_id path string true "Movie ID"
// @Param movie body models.Movie true "Movie object"
// @Success 200 {object} models.Movie
// @Failure 400 {object} map[string]interface{}
// @Router /movie/{movie_id} [put]
func UpdateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// DeleteMovie godoc
// @Summary Delete a movie
// @Description Delete a movie by ID
// @Tags Admin
// @Param movie_id path string true "Movie ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /movie/{movie_id} [delete]
func DeleteMovie(c *gin.Context) {
	// code to delete movie by ID
	c.Status(http.StatusNoContent)
}

// CreateMovieSchedule godoc
// @Summary Create a movie schedule
// @Description Create a schedule for a movie
// @Tags Admin
// @Param movie_id path string true "Movie ID"
// @Param schedule_id path string true "Schedule ID"
// @Success 201 {object} models.Schedule
// @Failure 400 {object} map[string]interface{}
// @Router /movie/{movie_id}/schedule/{schedule_id} [post]
func CreateMovieSchedule(c *gin.Context) {
	var schedule models.Schedule
	// code to create movie schedule
	c.JSON(http.StatusCreated, gin.H{"data": schedule})
}

// UpdateMovieSchedule godoc
// @Summary Update a movie schedule
// @Description Update a schedule for a movie
// @Tags Admin
// @Param movie_id path string true "Movie ID"
// @Param schedule_id path string true "Schedule ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /movie/{movie_id}/schedule/{schedule_id} [put]
func UpdateMovieSchedule(c *gin.Context) {
	// code to update movie schedule
	c.Status(http.StatusNoContent)
}

// DeleteSchedule godoc
// @Summary Delete movie schedule
// @Description Delete a specific movie schedule by id
// @Tags Admin
// @Param movie_id path string true "Movie ID"
// @Param schedule_id path string true "Schedule ID"
// @Produce json
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Router /movie/{movie_id}/schedule/{schedule_id} [delete]
func DeleteSchedule(c *gin.Context) {
	// code to delete schedule
	c.Status(http.StatusNoContent)
}

// CreateBranch godoc
// @Summary Create a new branch
// @Description Create a new branch with the given details
// @Tags Admin
// @Accept json
// @Produce json
// @Param branch body models.Branch true "Branch details"
// @Success 201 {object} models.Branch
// @Failure 400 {object} map[string]interface{}
// @Router /branch [post]
func CreateBranch(c *gin.Context) {
	var branch models.Branch
	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// code to create a new branch
	c.JSON(http.StatusCreated, gin.H{"data": branch})
}

// UpdateBranch godoc
// @Summary Update an existing branch
// @Description Update an existing branch with the given details
// @Tags Admin
// @Accept json
// @Produce json
// @Param branch_id path string true "Branch ID"
// @Param branch body models.Branch true "Branch details"
// @Success 200 {object} models.Branch
// @Failure 400 {object} map[string]interface{}
// @Router /branch/{branch_id} [put]
func UpdateBranch(c *gin.Context) {
	var branch models.Branch
	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// code to update the branch with branchID
	c.JSON(http.StatusOK, gin.H{"data": branch})
}

// DeleteBranch godoc
// @Summary Delete Branch
// @Description Delete Branch by ID
// @Tags Admin
// @Param branch_id path int true "Branch ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /branch/{branch_id} [delete]
func DeleteBranch(c *gin.Context) {
	// TODO: Implement branch deletion logic here
	c.Status(http.StatusNoContent)
}

// CreateTheatre godoc
// @Summary Create Theatre
// @Description Create a new Theatre in a Branch
// @Tags Admin
// @Accept json
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Param theatre body models.Theatre true "Theatre object"
// @Success 201 {object} models.Theatre
// @Failure 400 {object} map[string]interface{}
// @Router /branch/{branch_id}/theatre [post]
func CreateTheatre(c *gin.Context) {
	var theatre models.Theatre
	if err := c.ShouldBindJSON(&theatre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Implement theatre creation logic here
	c.JSON(http.StatusCreated, theatre)
}

// UpdateTheatre godoc
// @Summary Update Theatre
// @Description Update a Theatre in a Branch by ID
// @Tags Admin
// @Accept json
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Param theatre_id path int true "Theatre ID"
// @Param theatre body models.Theatre true "Theatre object"
// @Success 200 {object} models.Theatre
// @Failure 400 {object} map[string]interface{}
// @Router /branch/{branch_id}/theatre/{theatre_id} [put]
func UpdateTheatre(c *gin.Context) {
	var theatre models.Theatre
	if err := c.ShouldBindJSON(&theatre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Implement theatre update logic here
	c.JSON(http.StatusOK, theatre)
}

// DeleteTheatre godoc
// @Summary Delete Theatre
// @Description Delete a Theatre in a Branch by ID
// @Tags Admin
// @Param branch_id path int true "Branch ID"
// @Param theatre_id path int true "Theatre ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /branch/{branch_id}/theatre/{theatre_id} [delete]
func DeleteTheatre(c *gin.Context) {
	// TODO: Implement theatre deletion logic here
	c.Status(http.StatusNoContent)
}
