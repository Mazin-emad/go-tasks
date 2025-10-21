package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser     string
	DBPassword string
	DBName     string
	DBAddress  string
	JWTExpiresIn int64
	JWTKey string
}

var ConfigAmigoo = initConfig()

func initConfig() Config {

	godotenv.Load()

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost" ),
		Port: getEnv("PORT", "3306"),
		DBUser: getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "10Lofy@10"),
		DBName: getEnv("DB_NAME", "todo"),
		DBAddress: fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		JWTExpiresIn: getEnvInt("JWT_EXPIRES_IN", 3600*24*7),
		JWTKey: getEnv("JWT_KEY", "Chaos Is A Ladder"),
	}
}

func getEnv(key, fullback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fullback
}

func getEnvInt(key string, fullback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64); if err == nil {
			return i
		} else {
			return fullback
		}
	}
	return fullback
}

