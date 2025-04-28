package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/marcetin/gocms/db"
	"github.com/marcetin/gocms/routes"
)

func main() {
	// Initialize the database
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// Initialize Fiber with HTML template engine
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Static files
	app.Static("/static", "./static")

	// Admin Dashboard routes
	admin := app.Group("/admin")
	admin.Get("/", routes.AdminDashboard)        // Admin dashboard
	admin.Get("/create", routes.CreatePostPage)  // Create post page
	admin.Post("/create", routes.CreatePost)     // Create post
	admin.Get("/edit/:id", routes.EditPostPage)  // Edit post page
	admin.Post("/edit/:id", routes.EditPost)     // Edit post
	admin.Post("/delete/:id", routes.DeletePost) // Delete post

	// Frontend routes
	app.Get("/", routes.FrontendIndex)        // List all posts
	app.Get("/post/:id", routes.FrontendPost) // View a single post

	// Start the server
	log.Println("Starting server on :3000")
	log.Fatal(app.Listen(":3000"))
}
