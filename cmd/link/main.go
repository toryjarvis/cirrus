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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	database := db.Connect(dsn)

	app := fiber.New()

	routes.Register(app, database, []byte(jwtSecret))

	log.Fatal(app.Listen(":3001"))
}
