// models/upload.go
package models

import "time"

type Upload struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Nosbg           string    `json:"nosbg" gorm:"size:20;not null;index"`
	FileName        string    `json:"file_name"`
	FilePath        string    `json:"file_path"`
	HasilOCR        string    `json:"hasil_ocr"`
	HasilBacaan     string    `json:"hasil_bacaan"`
	UploaderID      int       `json:"uploader_id" gorm:"index"`
	Status          string    `json:"status" gorm:"type:varchar(20);index"`
	UploadedAt      time.Time `gorm:"column:uploaded_at;autoCreateTime" json:"uploaded_at,omitempty"`
	LastValidatedAt time.Time `gorm:"column:last_validated_at;autoUpdateTime" json:"last_validated_at,omitempty"`
}
