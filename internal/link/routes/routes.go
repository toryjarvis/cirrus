package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/toryjarvis/cirrus/internal/link/handlers"
	"github.com/toryjarvis/cirrus/internal/link/middleware"
)

// Register routes for link services
func Register(app *fiber.App, db *sqlx.DB, jwtSecret []byte) {
	authHandler := &handlers.AuthHandler{DB: db, JWTSecret: jwtSecret}
	linkHandler := &handlers.LinkHandler{DB: db}
	workspaceHandler := &handlers.WorkspaceHandler{DB: db}

	//Auths
	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)

	//Links
	api := app.Group("/api", middleware.Protected(jwtSecret))
	//Post
	api.Post("/links", linkHandler.Create)
	//Get
	api.Get("/links", linkHandler.List)
	//Get by ID
	api.Get("/links/:id", linkHandler.Get)
	//Patch
	api.Patch("/links/:id", linkHandler.Update)
	//Delete
	api.Delete("/links/:id", linkHandler.Delete)

	//Workspaces
	//Post
	api.Post("/workspaces", workspaceHandler.Create)
	//Get
	api.Get("/workspaces", workspaceHandler.List)
	//Patch
	api.Patch("/workspaces/:id", workspaceHandler.Update)
	//Delete
	api.Delete("/workspaces/:id", workspaceHandler.Delete)
}
