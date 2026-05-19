package handlers
import (
	"os"
	"log"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/jmoiron/sqlx"
)

type AuthHandler struct {
	DB *sqlx.DB
}

type RegisterRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func getJWTSecret() []byte {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("JWT_SECRET not set")
    }
    return []byte(secret)
}

// Registration Handler
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not hash password"})
	}

	var id string
	err = h.DB.QueryRow(
		"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
		req.Email, string(hash),
	).Scan(&id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not create user"})
	}

	return c.Status(201).JSON(fiber.Map{"id": id})
}

// Login Handler
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	var passwordHash string
	var id string
	err:= h.DB.QueryRow(
		"SELECT id, password_hash FROM users WHERE email = $1",
		req.Email,
	).Scan(&id, &passwordHash)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid email or password"})
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// JWT token gen
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString(getJWTSecret())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not sign token"})
	}

	return c.JSON(fiber.Map{"token": signedToken})

}