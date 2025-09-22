package services

import (
	"github.com/golang-jwt/jwt/v5"

	"server_aquascan/config"
)

// GenerateJWT wrapper biar pemanggilan lebih rapi dari controller
func GenerateJWT(userID uint, email string, role string) (string, error) {
	return config.GenerateJWT(userID, email, role)
}

// ParseJWT wrapper juga kalau butuh validasi manual
func ParseJWT(tokenString string) (*jwt.Token, error) {
	return config.ParseJWT(tokenString)
}
