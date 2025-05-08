package httpCommentV1

import (
	applicationComment "DDD/application/post_comment"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, service *applicationComment.PostCommentService) {
	handler := &Handler{Service: service}
	postGroup := app.Group("/api/v1/posts")

	postGroup.Get("/:postId/comments", handler.Paginate)
	postGroup.Post("/:postId/comments", handler.CreatePostComment)

	commentGroup := app.Group("/api/v1/comments")
	commentGroup.Get("/:id", handler.FindComment)
}
