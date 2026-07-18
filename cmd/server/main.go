package main

import (
    "log"
    "os"

    "github.com/joho/godotenv"

    "github.com/zuudevs/saq-inventory-system-backend/internal/config"
)

func main() {

    godotenv.Load()

    db, err := config.NewDatabase(
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_NAME"),
    )

    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    log.Println("Connected to MySQL")
}