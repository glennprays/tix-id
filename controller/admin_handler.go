package controller

import (
	"log"
	"net/http"
	"os"
	"time"
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
		tokenString, err := middleware.CreateToken(uint(admin.ID), "admin", os.Getenv("JWT_KEY"), 3600)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}

		middleware.SetCookie(c, "jwt_token", tokenString, time.Second*3600)
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
