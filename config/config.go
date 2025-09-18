package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecretKey []byte

func init() {
    if err := godotenv.Load(".env.dev"); err != nil {
        log.Fatal("Error loading .env.dev file")
    }

    secret := os.Getenv("JWT_SECRET_KEY")
    if secret == "" {
        log.Fatal("JWT_SECRET_KEY is not set in the environment")
    }
    JwtSecretKey = []byte(secret)
}
