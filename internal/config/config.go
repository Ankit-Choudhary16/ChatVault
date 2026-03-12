package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL string
    Port        string
    PrivateKey  string
}

func Load() Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using environment variables")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    return Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        Port:        port,
        PrivateKey:  os.Getenv("PRIVATE_KEY"),
    }
}