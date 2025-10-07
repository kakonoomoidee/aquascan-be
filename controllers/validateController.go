package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"server_aquascan/config"
	"server_aquascan/models"
	"server_aquascan/utils"
)

// Handler GET - ambil semua upload yang masih "submitted"
func GetSubmittedUploads(c *gin.Context) {
	var uploads []models.Upload

	if err := config.DB.
		Where("status = ?", "submitted").
		Order("uploaded_at DESC").
		Find(&uploads).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch pending uploads", err.Error())
		return
	}

	utils.RespondSuccess(c, uploads, "Pending uploads fetched successfully")
}

// Handler POST - validasi upload
func ValidateUpload(c *gin.Context) {
	var req struct {
		UploadID          uint   `json:"upload_id"`
		IsValid           bool   `json:"is_valid"`
		ValidationMessage string `json:"validation_message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Ambil user_id dari context
	adminVal, exists := c.Get("user_id")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "User belum terautentikasi", nil)
		return
	}

	// Convert user_id ke int
	var adminID int
	switch v := adminVal.(type) {
	case int:
		adminID = v
	case float64:
		adminID = int(v)
	case string:
		idInt, err := strconv.Atoi(v)
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, "Invalid user_id format", err.Error())
			return
		}
		adminID = idInt
	default:
		utils.RespondError(c, http.StatusBadRequest, "Unknown user_id type", nil)
		return
	}

	// Pastikan upload exist
	var upload models.Upload
	if err := config.DB.First(&upload, req.UploadID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Upload data not found", err.Error())
		return
	}

	// Simpan ke tabel upload_validations
	validation := models.UploadValidation{
		UploadID:          req.UploadID,
		AdminID:           adminID,
		IsValid:           req.IsValid,
		ValidationMessage: req.ValidationMessage,
		ValidatedAt:       time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := config.DB.Create(&validation).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to save validation record", err.Error())
		return
	}

	// Update status upload
	status := "invalid"
	if req.IsValid {
		status = "valid"
	}

	if err := config.DB.Model(&models.Upload{}).
		Where("id = ?", req.UploadID).
		Updates(map[string]interface{}{
			"status":            status,
			"hasil_bacaan":      gorm.Expr("CASE WHEN ? THEN hasil_ocr ELSE hasil_bacaan END", req.IsValid),
			"last_validated_at": time.Now(),
		}).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to update upload status", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"upload_id": req.UploadID,
		"status":    status,
	}, "Upload validation saved successfully")
}

// Handler GET - hitung berapa banyak upload yang masih "submitted"
func CountSubmittedUploads(c *gin.Context) {
	var count int64

	if err := config.DB.
		Model(&models.Upload{}).
		Where("status = ?", "submitted").
		Count(&count).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to count submitted uploads", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"count": count,
	}, "Counted submitted uploads successfully")
}
