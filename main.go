package main

import (
	"kajianku_be/config"
	"kajianku_be/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.InitUploader()
	e := routes.New()
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
