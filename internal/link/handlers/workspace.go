package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type WorkspaceHandler struct {
	DB *sqlx.DB
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

type UpdateWorkspaceRequest struct {
	Name string `json:"name"`
}

// Create
func (h *WorkspaceHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req CreateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	var id string
	err := h.DB.QueryRow(`
	INSERT INTO workspaces (user_id, name)
	VALUES ($1, $2)
	RETURNING id`, userID, req.Name,
	).Scan(&id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not create workspace"})
	}

	return c.Status(201).JSON(fiber.Map{"id": id, "name": req.Name})
}

// List
func (h *WorkspaceHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	rows, err := h.DB.Query(`
	SELECT id, name, created_at
	FROM workspaces
	WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not fetch workspaces"})
	}
	defer rows.Close()

	var workspaces []fiber.Map
	for rows.Next() {
		var id, name string
		var createdAt time.Time

		if err := rows.Scan(&id, &name, &createdAt); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "could not parse out workspaces"})
		}

		workspaces = append(workspaces, fiber.Map{
			"id":         id,
			"name":       name,
			"created_at": createdAt,
		})
	}
	return c.JSON(workspaces)
}

// Update
func (h *WorkspaceHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")

	var req UpdateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	_, err := h.DB.Exec(`
		UPDATE workspaces SET name = $1
		WHERE id = $2 AND user_id = $3`, req.Name, id, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not update workspace"})
	}

	return c.JSON(fiber.Map{"message": "workspace updated"})
}

// Delete
func (h *WorkspaceHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")

	result, err := h.DB.Exec(`DELETE FROM workspaces WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not delete workspace"})
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "workspace not found"})
	}

	return c.JSON(fiber.Map{"message": "workspace deleted"})
}
