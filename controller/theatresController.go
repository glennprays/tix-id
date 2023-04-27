package controller

import (
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ID of inserted branch"})
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
	// branchId := c.Query("branchId")
	// theatreId := c.Query("theatreId")
	var theatre models.Theatre
	if err := c.ShouldBindJSON(&theatre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var branch models.Branch
	responseData := models.BranchResponse{
		Response: models.Response{
			Status:  200,
			Message: "Theatre updated successfully",
		},
		Branch: branch,
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
	// branchId := c.Query("branchId")
	// theatreId := c.Query("theatreId")
	responseData := models.Response{
		Status:  200,
		Message: "Theatre deleted successfully",
	}
	c.JSON(http.StatusOK, responseData)
}
