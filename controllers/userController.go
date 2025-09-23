package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"server_aquascan/config"
	"server_aquascan/models"
	"server_aquascan/utils"
)

// GET ALL USERS (hanya admin)
func GetAllUsersHandler(c *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal mengambil data user", err.Error())
		return
	}

	utils.RespondSuccess(c, users, "Berhasil mengambil semua user")
}

// DELETE USER BY ID (soft delete, hanya admin)
func DeleteUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "ID user tidak valid", idParam)
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "User tidak ditemukan", err.Error())
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal menghapus user", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"id": id}, "User berhasil dihapus (soft delete)")
}

// UPDATE USER BY ID (hanya admin)
func UpdateUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "ID user tidak valid", idParam)
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "User tidak ditemukan", err.Error())
		return
	}

	var req struct {
		FullName string `json:"fullname"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Data request tidak valid", err.Error())
		return
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Gagal mengenkripsi password", err.Error())
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal update user", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"fullname": user.FullName,
		"role":     user.Role,
	}, "User berhasil diupdate")
}
