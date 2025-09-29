package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"server_aquascan/config"
	"server_aquascan/models"
	"server_aquascan/services"
	"server_aquascan/utils"
)

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func LoginHandler(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	var user models.User
	// Ambil kolom penting saja
	if err := config.DB.
		Select("id, password, email, role").
		Where("email = ?", payload.Email).
		First(&user).Error; err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Email atau password salah", nil)
		return
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Email atau password salah", nil)
		return
	}

	// Generate token pakai email & role juga
	token, err := services.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuat token", err.Error())
		return
	}

	// Return token + basic profile
	response := gin.H{
		"token": token,
		"user": gin.H{
			"id":   user.ID,
			"role": user.Role,
		},
	}

	utils.RespondSuccess(c, response, "Login berhasil")
}
