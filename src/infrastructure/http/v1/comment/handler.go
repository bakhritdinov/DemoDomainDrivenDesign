package httpCommentV1

import (
	applicationComment "DDD/src/application/post_comment"
	"DDD/src/domain"
	"DDD/src/domain/value_object"
	"DDD/src/infrastructure/http"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
	"strconv"
	"time"
)

type CreatePostCommentRequest struct {
	Text string `json:"text" example:"Great post"`
}

type PostCommentResponse struct {
	ID        uint       `json:"id" example:"1"`
	Text      string     `json:"text" example:"Great post"`
	PostId    uint       `json:"postId" example:"1"`
	CreatedAt time.Time  `json:"createdAt" swaggertype:"string" format:"date-time"`
	UpdatedAt time.Time  `json:"updatedAt" swaggertype:"string" format:"date-time"`
	DeletedAt *time.Time `json:"deletedAt" swaggertype:"string" format:"date-time"`
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
// @Param id path int true "post comment id"
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
		Text:      comment.Text.Value,
		PostId:    comment.PostId,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		DeletedAt: comment.DeletedAt,
	})
}

// Paginate paginate
// @Summary post comments pagination
// @Description post comments pagination
// @Tags posts
// @Accept json
// @Produce json
// @Param postId path int true "post id"
// @Param page query int false "page number" default(1)
// @Param per_page query int false "per page number" default(10)
// @Success 200 {object} http.PaginateResponse[domain.PostComment]
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

	return c.JSON(http.PaginateResponse[domain.PostComment]{
		Data: result.Comments,
		Pagination: http.Pagination{
			Page:       page,
			PerPage:    perPage,
			TotalItems: result.TotalCount,
			TotalPages: int(math.Ceil(float64(result.TotalCount) / float64(result.PerPage))),
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
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	postCommentText, err := value_object.NewPostCommentText(req.Text)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	comment, err := h.Service.CreatePostComment(c.Context(), domain.PostComment{
		Text:   postCommentText,
		PostId: uint(postId),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(PostCommentResponse{
		ID:        comment.Id,
		Text:      comment.Text.Value,
		PostId:    comment.PostId,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		DeletedAt: comment.DeletedAt,
	})
}
