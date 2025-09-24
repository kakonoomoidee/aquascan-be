package models

import "time"

type Client struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	BulanRek  string    `json:"bulan_rek"`
	Bulan     string    `json:"bulan"`
	TglBaca   time.Time `json:"tgl_baca"`
	Nosbg     string    `json:"nosbg" gorm:"size:12;not null;index"`
	IdTarip   int       `json:"idtarip"`
	Nama      string    `json:"nama"`
	Alamat    string    `json:"alamat"`
	Long      string    `json:"long"`
	Lat       string    `json:"lat"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Client) TableName() string {
	return "mas_bacahp"
}
