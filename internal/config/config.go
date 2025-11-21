package config

import (
    "log"

    "github.com/joho/godotenv"
    "os"
)

type Config struct {
    MongoURI  string
    JWTSecret string
    Port      string
}

var AppConfig Config

func LoadConfig() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    AppConfig = Config{
        MongoURI:  os.Getenv("MONGODB_URI"),
        JWTSecret: os.Getenv("JWT_SECRET"),
        Port:      os.Getenv("PORT"),
    }

    if AppConfig.Port == "" {
        AppConfig.Port = "8080"
    }
    if AppConfig.JWTSecret == "" {
        log.Fatal("JWT_SECRET required")
    }
    if AppConfig.MongoURI == "" {
        log.Fatal("MONGODB_URI required")
    }
}