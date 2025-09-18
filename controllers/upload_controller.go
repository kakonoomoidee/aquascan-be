package controllers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"server_aquascan/utils"
)

func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "File tidak ditemukan", err.Error())
		return
	}

	uploadDir := "uploads/temp"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuat folder upload", err.Error())
		return
	}

	newFileName := uuid.New().String() + filepath.Ext(file.Filename)
	savePath := filepath.Join(uploadDir, newFileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal menyimpan file", err.Error())
		return
	}

	userID, _ := c.Get("user_id")

	utils.RespondSuccess(c, gin.H{
		"user_id":   userID,
		"file_name": newFileName,
		"path":      savePath,
	}, "Upload berhasil")
}
