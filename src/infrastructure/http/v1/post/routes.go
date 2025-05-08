package httpPostV1

import (
	"DDD/src/application/post"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, service *applicationPost.PostService) {
	handler := &Handler{Service: service}
	postGroup := app.Group("/api/v1/posts")

	postGroup.Get("/", handler.Paginate)
	postGroup.Get("/:id", handler.FindPost)
	postGroup.Post("/", handler.CreatePost)
	postGroup.Patch("/:id", handler.UpdatePost)
	postGroup.Delete("/:id", handler.DeletePost)

}
