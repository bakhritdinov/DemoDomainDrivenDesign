package main

import (
	appPost "DDD/src/application/post"
	appComment "DDD/src/application/post_comment"
	"DDD/src/infrastructure/http/v1/comment"
	"DDD/src/infrastructure/http/v1/post"
	"DDD/src/infrastructure/persistence"
	"DDD/src/infrastructure/repository"
	"fmt"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/joho/godotenv"
	"os"
)

// @title Fiber Swagger
// @version 1.0
// @description Demo Domain driven design project.
// @BasePath /api
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})

	db, err := persistence.NewGormConnection(os.Getenv("DB_CONNECTION"))
	if err != nil {
		panic(err)
	}

	// Compress
	app.Use(compress.New())

	// E-tag
	app.Use(etag.New(etag.Config{
		Weak: true,
	}))

	// Health check
	app.Use(healthcheck.New())

	// Services
	postService := &appPost.PostService{
		PostRepo: repository.NewPostRepository(db),
	}
	commentService := &appComment.PostCommentService{
		PostCommentRepo: repository.NewCommentRepository(db),
		PostRepo:        repository.NewPostRepository(db),
	}

	// V1: Routes
	httpPostV1.SetupRoutes(app, postService)
	httpCommentV1.SetupRoutes(app, commentService)

	// Init dev tools
	initDevTools(app)

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}

func initDevTools(app *fiber.App) {
	if os.Getenv("APP_ENV") == "development" {
		// Register Swagger route
		app.Get("/docs/*", swagger.New(swagger.Config{
			BasePath: "/",
			FilePath: "./docs/swagger.json",
			Path:     "/docs",
			Title:    os.Getenv("APP_NAME"),
		}))

		// Metrics
		app.Get("/metrics", monitor.New())

		// Get all registered routes
		routes := app.GetRoutes()

		// Print all routes
		for _, route := range routes {
			fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
		}
	}
}
