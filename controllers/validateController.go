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

// Handler GET - ambil semua upload yang masih "submitted" dengan pagination
func GetSubmittedUploads(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	nosbgFilter := c.Query("nosbg")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	const maxLimit = 100
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := (page - 1) * limit

	baseQuery := config.DB.Model(&models.Upload{}).Where("status = ?", "submitted")
	if nosbgFilter != "" {
		baseQuery = baseQuery.Where("nosbg = ?", nosbgFilter)
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to count submitted uploads", err.Error())
		return
	}

	var uploads []models.Upload
	if err := baseQuery.
		Order("uploaded_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&uploads).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch submitted uploads", err.Error())
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	utils.RespondSuccess(c, gin.H{
		"data": uploads,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}, "Submitted uploads fetched successfully")
}

// Handler GET - ambil detail upload dan data validasinya
func GetUploadValidationDetail(c *gin.Context) {
	// ambil param id dari URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.RespondError(c, http.StatusBadRequest, "Invalid upload ID", nil)
		return
	}

	// ambil upload-nya
	var upload models.Upload
	if err := config.DB.First(&upload, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, "Upload not found", nil)
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch upload", err.Error())
		return
	}

	// ambil riwayat validasinya (jika ada)
	var validations []models.UploadValidation
	if err := config.DB.
		Where("upload_id = ?", id).
		Order("validated_at DESC").
		Find(&validations).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch validation records", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"upload":      upload,
		"validations": validations,
	}, "Upload validation detail fetched successfully")
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
