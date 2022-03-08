package env

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

const (
	DbHost = "DB_HOST"
	DbPort = "DB_PORT"
	DbName = "DB_NAME"
	DbUser = "DB_USER"
	DbPass = "DB_PASS"
)

const BindAddress = "BIND_ADDRESS"
const Secret = "SECRET"

// Get an environment variable by name This should be used over os.Getenv to ensure that the dotenv file has been loaded
func Get(key string) string {
	return os.Getenv(key)
}

// GetFallback gets an envar but returns the fallback value if the expected one is not found
func GetFallback(key, fallback string) string {
	val := Get(key)

	if val == "" {
		return fallback
	}

	return val
}
