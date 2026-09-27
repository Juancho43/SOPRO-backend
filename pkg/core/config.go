package core

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DatabaseURL     string
	Environment     string
	FocusmateAPIKEY string
	Timezone        string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se encontró el archivo .env, leyendo variables directamente del sistema.")
	}
	return &Config{
		Port:            os.Getenv("PORT"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		Environment:     os.Getenv("ENVIRONMENT"),
		FocusmateAPIKEY: os.Getenv("FOCUSMATE_APIKEY"),
		Timezone:        os.Getenv("TIMEZONE"),
	}
}
