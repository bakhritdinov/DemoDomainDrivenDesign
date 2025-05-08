package main

import (
	appPost "DDD/application/post"
	appComment "DDD/application/post_comment"
	httpCommentV1 "DDD/infrastructure/http/v1/comment"
	httpPostV1 "DDD/infrastructure/http/v1/post"
	"DDD/infrastructure/persistence"
	"DDD/infrastructure/persistence/validators"
	"DDD/infrastructure/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/joho/godotenv"
	"os"
)

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

	// Metrics
	app.Get("/metrics", monitor.New())

	// Validators
	validators.InitValidator("en")
	validators.RegisterValidationCallbacks(db)

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

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
