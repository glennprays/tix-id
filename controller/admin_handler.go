package controller

import (
	"log"
	"net/http"
	"tix-id/config"
	"tix-id/middleware"
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
	db := config.ConnectDB()
	defer db.Close()

	var login models.LoginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := db.QueryRow("select id, username, name, email, phone, nik from admin where email = ? and password = ?",
		login.Email,
		login.Password)

	var admin models.Admin
	if err := row.Scan(&admin.ID, &admin.Username, &admin.Name, &admin.Email, &admin.Phone, &admin.NIK); err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	} else {
		middleware.CreateToken(c, uint(admin.ID), "admin", 3600)

		responseData := models.AdminResponse{
			Response: models.Response{
				Status:  200,
				Message: "Login successful",
			},
			Admin: admin,
		}

		c.JSON(http.StatusCreated, responseData)
	}

}

// LogoutAccount godoc
// @Summary Logout Account
// @Description Logout Account admin and customer
// @Tags Auth
// @Success 200 {string} string "{"message": "Logout successful"}"
// @Router /auth/logout [post]
func LogoutAccount(c *gin.Context) {
	middleware.ResetUserToken(c)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
