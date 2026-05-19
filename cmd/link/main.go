package main
import (
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/toryjarvis/cirrus/internal/link/db"
	"github.com/toryjarvis/cirrus/internal/link/handlers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://cirrus:password@127.0.0.1:5433/cirrus?sslmode=disable"
	}

	database := db.Connect(dsn)
	authHandler := &handlers.AuthHandler{DB: database}

	app := fiber.New()

	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)

	log.Fatal(app.Listen(":3001"))
}