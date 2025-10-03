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

type ClientDetail struct {
	ID          uint64     `json:"id" gorm:"primaryKey"`
	BulanRek    string     `json:"bulan_rek"`
	Bulan       string     `json:"bulan"`
	TglBaca     *time.Time `json:"tgl_baca"`
	Nosbg       string     `json:"nosbg" gorm:"size:12;not null;index"`
	IdTarip     int8       `json:"idtarip"`
	Nama        string     `json:"nama"`
	Alamat      string     `json:"alamat"`
	Gol         string     `json:"gol"`
	Gol1        string     `json:"gol1"`
	UkMeter     string     `json:"uk_meter"`
	UkMeter1    string     `json:"uk_meter1"`
	JumlahRp    float64    `json:"jumlahrp"`
	Retribusi   float64    `json:"retribusi"`
	DanaMtr     float64    `json:"danamtr"`
	BiayaAdm    float64    `json:"biayaadm"`
	Meterai     float64    `json:"meterai"`
	StanLalu    float64    `json:"stan_lalu"`
	Pembaca     string     `json:"pembaca"`
	Rayon       string     `json:"rayon"`
	RayonKet    string     `json:"rayonket"`
	Kel         int        `json:"kel"`
	MeterBaca   float64    `json:"meter_baca"`
	Pakai       float64    `json:"pakai"`
	Waktu       *time.Time `json:"waktu"`
	WaktuProses *time.Time `json:"waktu_proses"`
	Tafsir      string     `json:"tafsir"`
	Ket         string     `json:"ket"`
	Long        string     `json:"long"`
	Lat         string     `json:"lat"`
	Rata23Bln   float64    `json:"rata2_3bln"`
	TglDownload *time.Time `json:"tgl_download"`
	NoMeter     string     `json:"nometer"`
	Long1       string     `json:"long1"`
	Lat1        string     `json:"lat1"`
	RATA2       float64    `json:"rata2"`
}

func (Client) TableName() string {
	return "mas_bacahp"
}

func (ClientDetail) TableName() string {
	return "mas_bacahp"
}
