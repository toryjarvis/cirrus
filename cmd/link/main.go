package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/toryjarvis/cirrus/internal/link/db"
	"github.com/toryjarvis/cirrus/internal/link/routes"
)

func main() {
	godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	database := db.Connect(dsn)

	app := fiber.New()

	routes.Register(app, database)

	log.Fatal(app.Listen(":3001"))
}
