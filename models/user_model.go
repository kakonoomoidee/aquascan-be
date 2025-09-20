// models/user_model.go
package models

import (
	"time"

	"gorm.io/gorm"

	"server_aquascan/config"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FullName  string         `gorm:"size:100;not null" json:"fullname"`
	Email     string         `gorm:"size:100;unique;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // jangan diexpose
	Role      string         `gorm:"size:50;default:'user'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func MigrateUser() {
	config.DB.AutoMigrate(&User{})
}
