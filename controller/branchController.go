package controller

import (
	"net/http"
	"tix-id/models"

	"github.com/gin-gonic/gin"
)

// CreateBranch godoc
// @Summary Create a new branch
// @Description Create a new branch with the given details
// @Tags Admin/Branch
// @Accept json
// @Produce json
// @Param branch body models.Branch true "Branch details"
// @Success 201 {object} models.BranchResponse
// @Router /branch [post]
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
// @Tags Admin/Branch
// @Accept json
// @Produce json
// @Success 200 {object} models.BranchesResponse
// @Router /branches [get]
func GetBranches(c *gin.Context) {
	// branchId := c.Query("branchId")
	var branches []models.Branch
	responseData := models.BranchesResponse{
		Response: models.Response{
			Status:  200,
			Message: "Branch updated successfully",
		},
		Branches: branches,
	}
	c.JSON(http.StatusOK, responseData)
}

// GetBranche godoc
// @Summary Get branche
// @Description Get a branche by movie_id and branche_id.
// @Tags Admin/Branch
// @Param branchId path string true "branch id"
// @Accept json
// @Produce json
// @Success 200 {object} models.BranchesResponse
// @Router /branches/{branchId} [get]
func GetBranche(c *gin.Context) {
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
// @Tags Admin/Branch
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
// @Tags Admin/Branch
// @Accept json
// @Produce json
// @Param branchId path string true "Branch ID"
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
