package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"server_aquascan/config"
	"server_aquascan/models"
	"server_aquascan/utils"
)

func ProfileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "User tidak ditemukan dalam token", nil)
		return
	}

	var id uint
	switch v := userID.(type) {
	case float64:
		id = uint(v) 
	case int:
		id = uint(v)
	default:
		utils.RespondError(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "User tidak ditemukan", err.Error())
		return
	}

	// Jangan return password hash!
	utils.RespondSuccess(c, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"fullname": user.FullName,
		"role":     user.Role,
	}, "Profile berhasil diambil")
}
