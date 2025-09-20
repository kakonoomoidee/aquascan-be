// controllers/auth_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"server_aquascan/config"
	model "server_aquascan/models"
	"server_aquascan/utils"
)

type LoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterHandler(c *gin.Context) {
	var payload model.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	payload.Password = string(hashedPassword)

	if err := config.DB.Create(&payload).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuat user", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"user": payload}, "Registrasi berhasil")
}

func LoginHandler(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	var user model.User
	if err := config.DB.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Email atau password salah", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Email atau password salah", nil)
		return
	}

	token, err := config.GenerateJWT(user.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuat token", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"token": token}, "Login berhasil")
}
