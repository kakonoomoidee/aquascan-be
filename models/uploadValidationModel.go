package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UploadValidation struct {
	ID                uint   `gorm:"primaryKey"`
	UploadID          uint   `json:"upload_id"`
	AdminID           int    `json:"admin_id"`
	IsValid           bool   `json:"is_valid"`
	ValidationMessage string `json:"validation_message"`
	ValidatedAt       string `gorm:"column:validated_at;autoUpdateTime" json:"validated_at,omitempty"`
}

func (UploadValidation) TableName() string {
	return "upload_validations"
}

func (v *UploadValidation) AfterCreate(tx *gorm.DB) (err error) {
	if v.IsValid {
		var upload Upload
		if err := tx.First(&upload, v.UploadID).Error; err != nil {
			fmt.Println("[Hook] Upload tidak ditemukan:", err)
			return err
		}

		// ✅ Pastikan hasil bacaan dikonversi ke float64 valid
		valStr := strings.ReplaceAll(upload.HasilBacaan, ",", ".")
		meterBaca, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			fmt.Println("[Hook] Format hasil_bacaan tidak valid:", upload.HasilBacaan)
			return fmt.Errorf("invalid numeric format for hasil_bacaan: %s", upload.HasilBacaan)
		}

		updateData := map[string]interface{}{
			"meter_baca":   meterBaca,
			"waktu_proses": time.Now(),
		}

		if err := tx.Table("mas_bacahp").
			Where("nosbg = ?", upload.Nosbg).
			Updates(updateData).Error; err != nil {
			fmt.Println("[Hook] Gagal update mas_bacahp:", err)
			return err
		}

		fmt.Printf("[Hook] ✅ Berhasil update mas_bacahp nosbg=%s meter_baca=%.2f\n", upload.Nosbg, meterBaca)
	}

	return nil
}
