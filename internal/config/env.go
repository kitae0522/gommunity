package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret              string
	JWTExpirationInSeconds int64
	RedisHost              string
	RedisPassword          string
	RedisDB                int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		JWTSecret:              getEnv("JWT_SECRET", "tempSecret"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 60*60*24*7),
		RedisHost:              getEnv("REDIS_HOST", "tempHost"),
		RedisPassword:          getEnv("REDIS_PASSWORD", "tempPassword"),
		RedisDB:                getEnvAsInt("REDIS_DB", 0),
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
