package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"server_aquascan/config"
	"server_aquascan/models"
	"server_aquascan/utils"
)

type Detection struct {
	Text string `json:"text"`
}

type OCRResponse struct {
	FileName   string      `json:"file_name"`
	Detections []Detection `json:"detections"`
}

// user mengupload gambar untuk diproses OCR
func OCRHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "File tidak ditemukan", err.Error())
		return
	}

	nosbg := c.PostForm("nosbg")
	if nosbg == "" {
		utils.RespondError(c, http.StatusBadRequest, "nosbg tidak boleh kosong", nil)
		return
	}

	// ambil user id dari JWT context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "User belum terautentikasi", nil)
		return
	}

	var userID int
	switch v := userIDValue.(type) {
	case string:
		userID, err = strconv.Atoi(v)
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Gagal konversi user_id ke int", err.Error())
			return
		}
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		utils.RespondError(c, http.StatusInternalServerError, "Tipe data user_id tidak dikenali", nil)
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

	pythonOCRURL := os.Getenv("PYTHON_OCR_URL")
	if pythonOCRURL == "" {
		utils.RespondError(c, http.StatusInternalServerError, "Python OCR URL tidak ditemukan di env", nil)
		return
	}

	f, err := os.Open(savePath)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuka file", err.Error())
		return
	}
	defer f.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, _ := w.CreateFormFile("file", newFileName)
	if _, err := io.Copy(fw, f); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal copy file ke request OCR", err.Error())
		return
	}
	w.Close()

	resp, err := http.Post(pythonOCRURL+"/ocr", w.FormDataContentType(), &b)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal request ke OCR service", err.Error())
		return
	}
	defer resp.Body.Close()

	var ocrResp OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal decode respons OCR", err.Error())
		return
	}

	var hasilOCR string
	if len(ocrResp.Detections) > 0 {
		hasilOCR = ocrResp.Detections[0].Text
	}

	fileBaseURL := os.Getenv("FILE_BASE_URL")

	publicFileURL := fileBaseURL + "temp/" + newFileName

	upload := models.Upload{
		Nosbg:      nosbg,
		FileName:   newFileName,
		FilePath:   publicFileURL,
		HasilOCR:   hasilOCR,
		UploaderID: userID,
		Status:     "pending",
	}

	if err := config.DB.Create(&upload).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal menyimpan ke database", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"id":         upload.ID,
		"nosbg":      upload.Nosbg,
		"file_name":  upload.FileName,
		"file_url":   publicFileURL,
		"hasil_ocr":  hasilOCR,
		"status":     upload.Status,
		"uploaderID": upload.UploaderID,
	}, "OCR berhasil diproses dan disimpan")
}

// user mengirim hasil bacaan yang sudah dikonfirmasi
func SubmitOCRHandler(c *gin.Context) {
	var req struct {
		Nosbg       string `json:"nosbg"`
		HasilBacaan string `json:"hasil_bacaan"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// cari data upload berdasarkan nosbg
	var upload models.Upload
	if err := config.DB.Where("nosbg = ?", req.Nosbg).First(&upload).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Upload data not found for given nosbg", err.Error())
		return
	}

	// update hasil bacaan dan status
	upload.HasilBacaan = req.HasilBacaan
	upload.Status = "submitted"

	if err := config.DB.Save(&upload).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to save confirmed result", err.Error())
		return
	}

	utils.RespondSuccess(c, upload, "OCR result successfully confirmed and saved")
}
