package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"server_aquascan/config"
	"server_aquascan/models"
	"server_aquascan/utils"
)

// 1️⃣ SUBMITTED UPLOADS — refresh tiap 1 menit
func GetSubmittedUploadsCount(c *gin.Context) {
	var count int64
	if err := config.DB.
		Model(&models.Upload{}).
		Where("status = ?", "submitted").
		Count(&count).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to count submitted uploads", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"submitted_uploads": count}, "Submitted uploads count fetched successfully")
}

// 2️⃣ VALIDATED TODAY — refresh tiap 30 menit
func GetValidatedTodayCount(c *gin.Context) {
	today := time.Now().Format("2006-01-02")

	var count int64
	if err := config.DB.
		Model(&models.UploadValidation{}).
		Where("DATE(validated_at) = ?", today).
		Count(&count).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to count validated uploads today", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"validated_today": count}, "Validated uploads today fetched successfully")
}

// 3️⃣ ACTIVE OFFICERS — refresh tiap 2 jam
func GetActiveOfficersCount(c *gin.Context) {
	var staffCount int64
	var adminCount int64

	// Hitung petugas (staff)
	if err := config.DB.
		Model(&models.User{}).
		Where("role = ?", "staff").
		Count(&staffCount).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to count staff", err.Error())
		return
	}

	// Hitung admin (opsional)
	if err := config.DB.
		Model(&models.User{}).
		Where("role = ?", "admin").
		Count(&adminCount).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to count admin", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"officers_active": staffCount, // hanya staff aktif
		"admin_total":     adminCount, // opsional, bisa diabaikan di FE
	}, "Active officers count fetched successfully")
}

// 4️⃣ TOTAL SUBMISSION — refresh tiap 1 jam
func GetTotalSubmissionsCount(c *gin.Context) {
	var count int64
	if err := config.DB.
		Model(&models.Upload{}).
		Count(&count).Error; err != nil && err != gorm.ErrRecordNotFound {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to count total submissions", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"total_submission": count}, "Total submission count fetched successfully")
}
