package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiKey      string
	OpenWeatherKey string
}

var App Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	App = Config{
		GeminiKey:      getEnv("GEMINI_API_KEY"),
		OpenWeatherKey: getEnv("OPENWEATHER_API_KEY"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env variable: %s", key)
	}
	return val
}
