package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"server_aquascan/config"
	"server_aquascan/models"
)

func GetClientsHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	search := c.Query("search")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// batasin limit cuma 10, 25, 50, 100
	if limit != 10 && limit != 25 && limit != 50 && limit != 100 {
		limit = 10
	}

	var clients []models.Client
	var totalItems int64

	query := config.DB.Model(&models.Client{})

	// kalau search ada, kasih WHERE
	if search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where("nosbg LIKE ? OR nama LIKE ? OR alamat LIKE ?", likeSearch, likeSearch, likeSearch)
	}

	// hitung total data setelah filter
	if err := query.Count(&totalItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hitung total data"})
		return
	}

	offset := (page - 1) * limit

	// ambil data dengan pagination
	if err := query.Select("id, nosbg, nama, alamat, `long`, lat").
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&clients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data"})
		return
	}

	totalPages := int((totalItems + int64(limit) - 1) / int64(limit)) // ceiling division

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"results":     clients,
			"totalPages":  totalPages,
			"totalItems":  totalItems,
			"currentPage": page,
		},
	})
}

// GET /api/clients/:nosbg
func GetClientDetailHandler(c *gin.Context) {
	nosbg := c.Param("nosbg")

	var client models.Client
	if err := config.DB.Where("nosbg = ?", nosbg).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": client,
	})
}

// GET /api/admin/clients/:nosbg
func GetMoreClientDetailHandler(c *gin.Context) {
	nosbg := c.Param("nosbg")

	var client models.ClientDetail
	if err := config.DB.Where("nosbg = ?", nosbg).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": client,
	})
}
