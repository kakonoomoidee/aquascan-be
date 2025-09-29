package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"server_aquascan/utils"
)

type Detection struct {
	Text string `json:"text"`
}

type OCRResponse struct {
	FileName   string      `json:"file_name"`
	Detections []Detection `json:"detections"`
}

func UploadHandler(c *gin.Context) {
	// ambil file dari form
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "File tidak ditemukan", err.Error())
		return
	}

	// bikin folder upload
	uploadDir := "uploads/temp"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal membuat folder upload", err.Error())
		return
	}

	// simpan dengan nama random
	newFileName := uuid.New().String() + filepath.Ext(file.Filename)
	savePath := filepath.Join(uploadDir, newFileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal menyimpan file", err.Error())
		return
	}

	// ambil user id dari context
	userID, _ := c.Get("user_id")

	// kirim file ke Python service (form-data)
	pythonOCRURL := os.Getenv("PYTHON_OCR_URL")
	if pythonOCRURL == "" {
		utils.RespondError(c, http.StatusInternalServerError, "Python OCR URL tidak ditemukan di env", nil)
		return
	}

	// buka file
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
		utils.RespondError(c, http.StatusInternalServerError, "Gagal copy file ke request", err.Error())
		return
	}
	w.Close()

	// request ke Python
	resp, err := http.Post(pythonOCRURL+"/ocr", w.FormDataContentType(), &b)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal request ke OCR service", err.Error())
		return
	}
	defer resp.Body.Close()

	// decode respons dari Python
	var ocrResp OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Gagal decode respons OCR", err.Error())
		return
	}

	// balikin ke user sesuai format standar
	utils.RespondSuccess(c, gin.H{
		"user_id":    userID,
		"file_name":  ocrResp.FileName,
		"detections": ocrResp.Detections,
	}, "OCR berhasil diproses")
}
