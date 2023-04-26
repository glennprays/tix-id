package controller

import (
	"fmt"
	"net/http"
	"tix-id/config"
	"tix-id/models"

	"github.com/gin-gonic/gin"
)

// CreateBranch godoc
// @Summary Create a new branch
// @Description Create a new branch with the given details
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body models.Branch true "Branch details"
// @Success 201 {object} models.BranchResponse
// @Router /branches [post]
func CreateBranch(c *gin.Context) {
	var branch models.Branch
	if err := c.ShouldBindJSON(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseData := models.BranchResponse{
		Response: models.Response{
			Status:  200,
			Message: "Branch inserted successfully",
		},
		Branch: branch,
	}
	c.JSON(http.StatusOK, responseData)
}

// GetBranches godoc
// @Summary Get branches
// @Description Get a branches by movie_id and branches_id.
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} models.BranchesResponse
// @Router /branches [get]
func GetBranches(c *gin.Context) {
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()

	// Execute a SELECT query to retrieve all branches from the database
	rows, err := db.Query("SELECT * FROM branch")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve branches from database"})
		return
	}

	// Iterate over the rows returned from the query and store them in a slice of Branch structs
	var branches []models.Branch
	for rows.Next() {
		var branch models.Branch
		err := rows.Scan(&branch.ID, &branch.Name, &branch.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve branches from database"})
			return
		}

		// Execute a SELECT query to retrieve all theatres for the current branch
		theatreRows, err := db.Query("SELECT id, name FROM theatre WHERE branch_id = ?", branch.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": err})
			return
		}

		// Create an empty slice of Theatre structs to hold the retrieved theatres
		var theatres []models.Theatre

		// Iterate over the theatre rows returned from the query and store them in the theatres slice
		for theatreRows.Next() {
			var theatre models.Theatre
			err := theatreRows.Scan(&theatre.ID, &theatre.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error2": fmt.Sprintf("%v", err)})
				return
			}
			theatres = append(theatres, theatre)
		}

		// Add the retrieved theatres to the current branch
		branch.Theatres = &theatres

		// Add the current branch to the branches slice
		branches = append(branches, branch)
	}

	// Create a BranchesResponse struct with the retrieved branches and send it as a JSON response
	responseData := models.BranchesResponse{
		Response: models.Response{
			Status:  200,
			Message: "Branches retrieved successfully",
		},
		Branches: branches,
	}
	c.JSON(http.StatusOK, responseData)
}

// GetBranch godoc
// @Summary Get branche
// @Description Get a branche by movie_id and branche_id.
// @Tags Admin
// @Param branchId path string true "branch id"
// @Accept json
// @Produce json
// @Success 200 {object} models.BranchesResponse
// @Router /branches/{branchId} [get]
func GetBranch(c *gin.Context) {
	// branchId := c.Query("branchId")
	var branch models.Branch
	responseData := models.BranchResponse{
		Response: models.Response{
			Status:  200,
			Message: "Branch updated successfully",
		},
		Branch: branch,
	}
	c.JSON(http.StatusOK, responseData)
}

// DeleteBranch godoc
// @Summary Delete Branch
// @Description Delete Branch by ID
// @Tags Admin
// @Param branchId path int true "Branch ID"
// @Success 204 {object} models.Response
// @Router /branches/{branchId} [delete]
func DeleteBranch(c *gin.Context) {
	// branchId := c.Query("branchId")
	responseData := models.Response{
		Status:  200,
		Message: "Movie deleted successfully",
	}
	c.JSON(http.StatusOK, responseData)
}

// UpdateBranch godoc
// @Summary Update an existing branch
// @Description Update an existing branch with the given details
// @Tags Admin
// @Accept json
// @Produce json
// @Param branchId path string true "Branch ID"
// @Param body body models.Branch true "Branch details"
// @Success 200 {object} models.BranchResponse
// @Router /branches/{branchId} [put]
func UpdateBranch(c *gin.Context) {
	// branchId := c.Query("branchId")
	var branch models.Branch
	responseData := models.BranchResponse{
		Response: models.Response{
			Status:  200,
			Message: "Branch updated successfully",
		},
		Branch: branch,
	}
	c.JSON(http.StatusOK, responseData)
}
