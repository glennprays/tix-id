package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"tix-id/config"
	"tix-id/models"

	"github.com/gin-gonic/gin"
)

// CreateTheatre godoc
// @Summary Create Theatre
// @Description Create a new Theatre in a Branch
// @Tags Admin
// @Accept json
// @Produce json
// @Param branchId path int true "Branch ID"
// @Param body body models.Theatre true "Theatre details"
// @Success 201 {object} models.BranchResponse
// @Router /branches/{branchId}/theatres [post]
func CreateTheatre(c *gin.Context) {
	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()

	// branchId := c.Query("branchId")
	branchId, err := strconv.Atoi(c.Param("branchId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve branches from database"})
		return
	}

	var theatre models.Theatre
	if errs := c.ShouldBindJSON(&theatre); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := db.Exec("INSERT INTO theatre (name,branch_id) VALUES (?,?)", theatre.Name, branchId)

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ID of inserted theatre"})
		return
	}
	theatre.ID = int(id)
	responseData := models.TheatreResponse{
		Response: models.Response{
			Status:  200,
			Message: "Theatre inserted successfully",
		},
		Theatre: theatre,
	}
	c.JSON(http.StatusOK, responseData)
}

// UpdateTheatre godoc
// @Summary Update Theatre
// @Description Update a Theatre in a Branch by ID
// @Tags Admin
// @Accept json
// @Produce json
// @Param branchId path int true "Branch ID"
// @Param theatreId path int true "Theatre ID"
// @Param body body models.Theatre true "Theatre details"
// @Success 200 {object} models.BranchResponse
// @Router /branches/{branchId}/theatres/{theatreId} [put]
func UpdateTheatre(c *gin.Context) {
	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()
	// branchId := c.Query("branchId")
	branchId, err := strconv.Atoi(c.Param("branchId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve branches from database"})
		return
	}
	// theatreId := c.Query("theatreId")
	theatreId, err := strconv.Atoi(c.Param("theatreId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve theatre from database"})
		return
	}
	var theatre models.Theatre
	if err := c.ShouldBindJSON(&theatre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var branch models.Branch
	//theatre id
	theatre.ID = theatreId
	//branch id
	branch.ID = branchId

	result, err := db.Exec("UPDATE theatre SET name=? WHERE id=? && branch_id=?", theatre.Name, theatre.ID, branch.ID)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Theatre not found"})
		return
	}

	responseData := models.TheatreResponse{
		Response: models.Response{
			Status:  200,
			Message: "Theatre updated successfully",
		},
		Theatre: theatre,
	}
	c.JSON(http.StatusOK, responseData)
}

// DeleteTheatre godoc
// @Summary Delete Theatre
// @Description Delete a Theatre in a Branch by ID
// @Tags Admin
// @Param branchId path int true "Branch ID"
// @Param theatreId path int true "Theatre ID"
// @Success 204 {object} models.Response
// @Router /branches/{branchId}/theatres/{theatreId} [delete]
func DeleteTheatre(c *gin.Context) {

	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()
	// branchId := c.Query("branchId")
	branchId, err := strconv.Atoi(c.Param("branchId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Theatre from database"})
		return
	}
	// theatreId := c.Query("theatreId")

	theatreId, err := strconv.Atoi(c.Param("theatreId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Theatre from database"})
		return
	}

	// Check if the branch exists
	var count int
	errs := db.QueryRow("SELECT COUNT(*) FROM branch WHERE id=?", branchId).Scan(&count)
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theatre not found"})
		return
	}

	// Check if the theatre exists
	var counts int
	error := db.QueryRow("SELECT COUNT(*) FROM theatre WHERE id=? && branch_id=?", theatreId, branchId).Scan(&counts)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if counts == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "This Theatre doesn't have any theatre yet"})
		return
	}

	// Delete the theatre
	result, err := db.Exec("DELETE FROM theatre WHERE id=? && branch_id=?", theatreId, branchId)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete theatre"})
		return
	}
	message := fmt.Sprintf("Theatre with id %d in branch with id %d was deleted successfully", theatreId, branchId)
	responseData := models.Response{
		Status:  200,
		Message: message,
	}
	c.JSON(http.StatusOK, responseData)
}
