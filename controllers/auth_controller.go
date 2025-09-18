package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"server_aquascan/services"
	"server_aquascan/utils"
)

type LoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Input tidak valid", err.Error())
		return
	}

	// Dummy user (sementara hardcoded, nanti bisa dari DB)
	userEmailFromDB := "kaisar@kerajaan.com"
	passwordHashFromDB := "$2a$10$ZzL.L93GD3xr1NVKFrD1jOta30NHteRl7hgECczlXPM0u0DojP2wS" // hash "password123"

	if !strings.EqualFold(payload.Email, userEmailFromDB) {
		utils.RespondError(c, http.StatusUnauthorized, "Email atau password salah", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHashFromDB), []byte(payload.Password)); err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Email atau password salah", nil)
		return
	}

	token, err := services.GenerateJWT(1)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuat token", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"token": token}, "Login berhasil")
}

func ProfileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mendapatkan info user", nil)
		return
	}

	utils.RespondSuccess(c, gin.H{
		"user_id": userID,
	}, "Profile berhasil diambil")
}
