package httpCommentV1

import (
	applicationComment "DDD/src/application/post_comment"
	"DDD/src/domain"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
	"strconv"
	"time"
)

type CreatePostCommentRequest struct {
	Text string `json:"text" validate:"required,min=3,max=100" example:"Great post"`
}

type PostCommentResponse struct {
	ID        uint      `json:"id" example:"1"`
	Text      string    `json:"text" example:"Great post"`
	PostId    uint      `json:"postId" example:"1"`
	CreatedAt time.Time `json:"createdAt" swaggertype:"string" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" swaggertype:"string" format:"date-time"`
}

type Handler struct {
	Service *applicationComment.PostCommentService
}

// FindComment Find post comment
// @Summary Find post comment by id
// @Description Find post comment
// @Tags comments
// @Accept json
// @Produce json
// @Success 201 {object} PostCommentResponse
// @Failure 400 {string} error
// @Router /v1/comments/{id} [get]
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

	return c.JSON(PostCommentResponse{
		ID:        comment.Id,
		Text:      comment.Text,
		PostId:    comment.PostId,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	})
}

// Paginate paginate
// @Summary post comments pagination
// @Description post comments pagination
// @Tags posts
// @Accept json
// @Produce json
// @Param page path int false "page, default 1"
// @Param per_page path int false "per_page, default 10"
// @Success 200 {object} applicationPostComment.PaginatedComments
// @Failure 400 {string} error
// @Router /v1/posts/{postId}/comments [get]
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

// CreatePostComment create a new post comment data
// @Summary Create a new post comment
// @Description Create post comment
// @Tags posts
// @Accept json
// @Produce json
// @Param request body CreatePostCommentRequest true "Post comment data to create"
// @Success 201 {object} PostCommentResponse
// @Failure 400 {string} error
// @Router /v1/posts/{postId}/comments [post]
func (h *Handler) CreatePostComment(c *fiber.Ctx) error {
	postId, err := c.ParamsInt("postId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post id"})
	}

	req := CreatePostCommentRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post comment data"})
	}

	comment, err := h.Service.CreatePostComment(c.Context(), domain.PostComment{
		Text:   req.Text,
		PostId: uint(postId),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(PostCommentResponse{
		ID:        comment.Id,
		Text:      comment.Text,
		PostId:    comment.PostId,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	})
}
