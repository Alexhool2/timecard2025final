package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost              string
	Port                    string
	DbUser                  string
	DbPassword              string
	DbAddress               string
	DbName                  string
	JWTExpirationsInSeconds int64
	JWTSecret               string
	// Novas variáveis para o novo banco de dados
	NewDbUser     string
	NewDbPassword string
	NewDbAddress  string
	NewDbName     string
}

var Envs = initConfig()

func initConfig() Config {
	wd, _ := os.Getwd()
	fmt.Printf("working directory: %s\n", wd)
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file %v\n", err)
	} else {
		fmt.Println(".env file load successfully")
	}

	return Config{
		PublicHost:              getEnv("PUBLIC_HOST", "https://localhost"),
		Port:                    getEnv("PORT", "8081"),
		DbUser:                  getEnv("DB_USER", "root"),
		DbPassword:              getEnv("DB_PASSWORD", "myPassword"),
		DbAddress:               fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DbName:                  getEnv("DB_NAME", "timeCard2025"),
		JWTSecret:               getEnv("JWT_SECRET", "not-that-secret"),
		JWTExpirationsInSeconds: getEnvAsInt("JWT_EXP", 3600*24*1),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
