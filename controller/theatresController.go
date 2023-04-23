package controller

import (
	"net/http"
	"tix-id/models"

	"github.com/gin-gonic/gin"
)

// CreateTheatre godoc
// @Summary Create Theatre
// @Description Create a new Theatre in a Branch
// @Tags Admin/Branch/Theatre
// @Accept json
// @Produce json
// @Param branchId path int true "Branch ID"
// @Success 201 {object} models.TheatreResponse
// @Router /branches/{branchId}/theatres [post]
func CreateTheatre(c *gin.Context) {
	// branchId := c.Query("branchId")
	var theatre models.Theatre
	if err := c.ShouldBindJSON(&theatre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var branch models.Branch
	responseData := models.BranchResponse{
		Response: models.Response{
			Status:  200,
			Message: "Theatre inserted successfully",
		},
		Branch: branch,
	}
	c.JSON(http.StatusOK, responseData)
}

// UpdateTheatre godoc
// @Summary Update Theatre
// @Description Update a Theatre in a Branch by ID
// @Tags Admin/Branch/Theatre
// @Accept json
// @Produce json
// @Param branchId path int true "Branch ID"
// @Param theatreId path int true "Theatre ID"
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
// @Tags Admin/Branch/Theatre
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
