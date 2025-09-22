package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"server_aquascan/config"
	model "server_aquascan/models"
	"server_aquascan/services"
	"server_aquascan/utils"
)

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func RegisterHandler(c *gin.Context) {
	var payload model.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengenkripsi password", err.Error())
		return
	}
	payload.Password = string(hashedPassword)

	// Simpan user
	if err := config.DB.Create(&payload).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuat user", err.Error())
		return
	}

	// Jangan return password ke client
	response := gin.H{
		"id":    payload.ID,
		"email": payload.Email,
		"role":  payload.Role,
	}

	utils.RespondSuccess(c, gin.H{"user": response}, "Registrasi berhasil")
}

func LoginHandler(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	var user model.User
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
