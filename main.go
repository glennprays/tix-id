package main

import (
	"log"
	"tix-id/docs"
	"tix-id/routes"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	docs.SwaggerInfo.Title = "TIX-ID API Documentation"
	docs.SwaggerInfo.Description = "Ticketing Application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := routes.SetupRouter()
	r.Run(":8080")
}
