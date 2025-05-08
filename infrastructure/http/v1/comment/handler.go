package httpCommentV1

import (
	applicationComment "DDD/application/post_comment"
	"DDD/domain"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
	"strconv"
)

type Handler struct {
	Service *applicationComment.PostCommentService
}

func (h *Handler) FindComment(c *fiber.Ctx) error {
	commentId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment id"})
	}

	comment, err := h.Service.FindById(c.Context(), commentId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("PostComment with id %d not found", commentId)})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(comment)
}

func (h *Handler) Paginate(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	postId, err := c.ParamsInt("postId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post id"})
	}

	result, err := h.Service.FindPaginatedComments(c.Context(), postId, page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(fiber.Map{
		"data": result.Comments,
		"pagination": fiber.Map{
			"page":        result.Page,
			"per_page":    result.PerPage,
			"total_items": result.TotalCount,
			"total_pages": int(math.Ceil(float64(result.TotalCount) / float64(result.PerPage))),
		},
	})
}

func (h *Handler) CreatePostComment(c *fiber.Ctx) error {
	postId, err := c.ParamsInt("postId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post id"})
	}

	commentData := new(domain.PostComment)
	if err := c.BodyParser(commentData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post comment data"})
	}
	commentData.PostId = uint(postId)

	comment, err := h.Service.CreatePostComment(c.Context(), *commentData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(comment)
}
