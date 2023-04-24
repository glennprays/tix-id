package controller

import (
	"net/http"
	"tix-id/models"

	"github.com/gin-gonic/gin"
)

// LoginAdmin godoc
// @Summary Login Admin
// @Description Login Admin Account
// @Tags Admin
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Login details"
// @Success 200 {object} models.AdminResponse
// @Router /admin/auth/login [post]
func LoginAdmin(c *gin.Context) {
	var login models.LoginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
	responseData := models.AdminResponse{
		Response: models.Response{
			Status:  200,
			Message: "Login successful",
		},
		Admin: admin,
	}

	c.JSON(http.StatusCreated, responseData)
}
