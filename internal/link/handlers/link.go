package handlers

import (
	"log"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type LinkHandler struct{ DB *sqlx.DB }

// create link request
type CreateLinkRequest struct {
	OriginalURL string  `json:"original_url"`
	CustomSlug  *string `json:"custom_slug"`
	WorkspaceID string  `json:"workspace_id"`
	ExpiresAt   *string `json:"expires_at"`
}

// update link request
type UpdateLinkRequest struct {
	ExpiresAt *string `json:"expires_at"`
	IsActive  *bool   `json:"is_active"`
}

// slug generation
const slugCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateSlug(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = slugCharset[rand.IntN(len(slugCharset))]
	}

	return string(b)
}

// Create
func (h *LinkHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req CreateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	var wsExists bool
	//log.Printf("checking workspace: id=%s user_id=%s", req.WorkspaceID, userID)
	err := h.DB.QueryRow(`
	SELECT EXISTS(SELECT 1 FROM workspaces WHERE id = $1 AND user_id = $2)`,
		req.WorkspaceID, userID,
	).Scan(&wsExists)
	//log.Printf("wsExists: %v err: %v", wsExists, err)
	if err != nil || !wsExists {
		return c.Status(404).JSON(fiber.Map{"error": "workspace not found"})
	}

	slug := generateSlug(7)
	if req.CustomSlug != nil && *req.CustomSlug != "" {
		slug = *req.CustomSlug
	}

	var id string
	for i := 0; i < 5; i++ {
		err = h.DB.QueryRow(`
		INSERT INTO links (workspace_id, user_id, original_url, slug, custom_slug, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
			req.WorkspaceID, userID, req.OriginalURL, slug, req.CustomSlug, req.ExpiresAt,
		).Scan(&id)
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "duplicate key") {
			return c.Status(500).JSON(fiber.Map{"error": "could not create link"})
		}
		if req.CustomSlug != nil && *req.CustomSlug != "" {
			return c.Status(409).JSON(fiber.Map{"error": "slug already taken"})
		}
		slug = generateSlug(7)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not generate unique slug"})
	}

	return c.Status(201).JSON(fiber.Map{"id": id, "slug": slug})
}

// List
func (h *LinkHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	rows, err := h.DB.Query(`
	SELECT id, workspace_id, original_url, slug, custom_slug, expires_at, is_active, created_at
	FROM links WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		log.Printf("query error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "could not fetch links"})
	}
	defer rows.Close()

	var links []fiber.Map
	for rows.Next() {
		var id, workspaceID, originalURL, slug string
		var customSlug *string
		var expiresAt *time.Time
		var isActive bool
		var createdAt time.Time

		if err := rows.Scan(&id, &workspaceID, &originalURL, &slug, &customSlug, &expiresAt, &isActive, &createdAt); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "could not parse out links"})
		}

		links = append(links, fiber.Map{
			"id":           id,
			"workspace_id": workspaceID,
			"original_url": originalURL,
			"slug":         slug,
			"custom_slug":  customSlug,
			"expires_at":   expiresAt,
			"is_active":    isActive,
			"created_at":   createdAt,
		})
	}

	return c.JSON(links)
}

// Get by ID
func (h *LinkHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")

	var originalURL, slug string
	var customSlug *string
	var expiresAt *time.Time
	var isActive bool
	var createdAt time.Time

	err := h.DB.QueryRow(`
	SELECT original_url, slug, custom_slug, expires_at, is_active, created_at
	FROM links WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&originalURL, &slug, &customSlug, &expiresAt, &isActive, &createdAt)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "link not found"})
	}

	return c.JSON(fiber.Map{
		"id":           id,
		"original_url": originalURL,
		"slug":         slug,
		"custom_slug":  customSlug,
		"expires_at":   expiresAt,
		"is_active":    isActive,
		"created_at":   createdAt,
	})
}

// Update
func (h *LinkHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")

	var req UpdateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	_, err := h.DB.Exec(`
		UPDATE links SET is_active = COALESCE($1, is_active), expires_at = COALESCE($2, expires_at)
		WHERE id = $3 AND user_id = $4`, req.IsActive, req.ExpiresAt, id, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not update link"})
	}

	return c.JSON(fiber.Map{"message": "link updated"})

}

// Delete
func (h *LinkHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")

	result, err := h.DB.Exec(`DELETE FROM links WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not delete link"})
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "link not found"})
	}

	return c.JSON(fiber.Map{"message": "link deleted"})
}
