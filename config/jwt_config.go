package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	JwtSecretKey   []byte
	JwtExpireHours int
)

func InitJWT() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found, pakai system env")
	}

	// Ambil secret
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("❌ JWT_SECRET_KEY tidak boleh kosong di .env")
	}
	JwtSecretKey = []byte(secret)

	// Ambil expire hours
	expireHourStr := os.Getenv("JWT_EXPIRE_HOURS")
	expireHour, err := strconv.Atoi(expireHourStr)
	if err != nil || expireHour <= 0 {
		log.Println("⚠️ JWT_EXPIRE_HOURS invalid atau kosong, fallback ke 24 jam")
		expireHour = 24
	}
	JwtExpireHours = expireHour
}

// GenerateJWT dengan userID, email, role
func GenerateJWT(userID uint, email string, role string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(JwtExpireHours))

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtSecretKey)
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JwtSecretKey, nil
	})
}
