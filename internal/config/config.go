package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Flags() (string, string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке .env файла")
	}
	address := flag.String("a", "localhost:8080", "адрес запуска HTTP-сервера")
	db := flag.String("d", "", "адрес для бд")
	api := flag.String("p", "", "адрес для api")
	flag.Parse()
	if envAddress := os.Getenv("SERVER_ADDRESS"); envAddress != "" {
		*address = envAddress
	}

	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		*db = envDB
	}

	if envApi := os.Getenv("API_ADDRESS"); envApi != "" {
		*api = envApi
	}

	return *address, *db, *api
}
