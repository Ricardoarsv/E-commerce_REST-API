package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost                 string
	DBPort                 string
	DBUser                 string
	DBPassword             string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	USERADMIN              int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		DBHost:                 getEnv("DB_HOST", "localhost"),
		DBPort:                 getEnv("DB_PORT", "5432"),
		DBUser:                 getEnv("DB_USER", "postgres"),
		DBPassword:             getEnv("DB_PASSWORD", "password"),
		DBName:                 getEnv("DB_NAME", "postgres"),
		JWTExpirationInSeconds: getEnvAtInt("JWTExpirationInSeconds", 3600*24*7),
		JWTSecret:              getEnv("JWTSecret", "secreto"),
		USERADMIN:              getEnvAtInt("USER_ADMIN", 1),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAtInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return i
		}
	}
	return fallback
}
