package httpPostV1

import (
	"DDD/src/application/post"
	"DDD/src/domain"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
	"strconv"
	"time"
)

type CreatePostRequest struct {
	Title   string `json:"title" example:"My Post Title"`
	Content string `json:"content" example:"Post content here"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" example:"My Post Title"`
	Content string `json:"content" example:"Post content here"`
}

type PostResponse struct {
	ID        uint      `json:"id" example:"1"`
	Title     string    `json:"title" example:"My Post Title"`
	Content   string    `json:"content" example:"Post content here"`
	CreatedAt time.Time `json:"createdAt" swaggertype:"string" format:"date-time"`
	UpdatedAt time.Time `json:"updatedAt" swaggertype:"string" format:"date-time"`
}

type Handler struct {
	Service *applicationPost.PostService
}

// FindPost Find post
// @Summary Find post by id
// @Description Find post
// @Tags posts
// @Accept json
// @Produce json
// @Success 201 {object} PostResponse
// @Failure 400 {string} error
// @Router /v1/posts/{id} [get]
func (h *Handler) FindPost(c *fiber.Ctx) error {
	postID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post id"})
	}

	post, err := h.Service.FindById(c.Context(), postID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Post with id %d not found", postID)})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(PostResponse{
		ID:        post.Id,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

// Paginate paginate
// @Summary posts pagination
// @Description posts pagination
// @Tags posts
// @Accept json
// @Produce json
// @Param page path int false "page, default 1"
// @Param per_page path int false "per_page, default 10"
// @Success 200 {object} applicationPost.PaginatedPosts
// @Failure 400 {string} error
// @Router /v1/posts [get]
func (h *Handler) Paginate(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	result, err := h.Service.FindPaginatedPosts(c.Context(), page, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(fiber.Map{
		"data": result.Posts,
		"pagination": fiber.Map{
			"page":        result.Page,
			"per_page":    result.PerPage,
			"total_items": result.TotalCount,
			"total_pages": int(math.Ceil(float64(result.TotalCount) / float64(result.PerPage))),
		},
	})
}

// CreatePost create a new post data
// @Summary Create a new post
// @Description Create post
// @Tags posts
// @Accept json
// @Produce json
// @Param request body CreatePostRequest true "Post data to create"
// @Success 201 {object} PostResponse
// @Failure 400 {string} error
// @Router /v1/posts [post]
func (h *Handler) CreatePost(c *fiber.Ctx) error {
	req := CreatePostRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	postData := domain.Post{
		Title:   req.Title,
		Content: req.Content,
	}

	post, err := h.Service.CreatePost(c.Context(), postData)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Post with title %s already exists", postData.Title)})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(PostResponse{
		ID:        post.Id,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

// UpdatePost update post
// @Summary Update post
// @Description Update post
// @Tags posts
// @Accept json
// @Produce json
// @Param request body UpdatePostRequest true "Post data to update"
// @Success 200 {object} PostResponse
// @Failure 400 {string} error
// @Router /v1/posts/{id} [patch]
func (h *Handler) UpdatePost(c *fiber.Ctx) error {
	postID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post id"})
	}
	post, err := h.Service.FindById(c.Context(), postID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Post with id %d not found", postID)})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	req := UpdatePostRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	post, err = h.Service.UpdatePost(c.Context(), domain.Post{
		Title:   req.Title,
		Content: req.Content,
	})

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You can't update a post with the same title."})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(PostResponse{
		ID:        post.Id,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

// DeletePost function removes a post by ID
// @Summary Remove post by ID
// @Description Remove post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 204 "No Content - Successful deletion"
// @Failure 400 {string} error
// @Router /v1/posts/{id} [delete]
func (h *Handler) DeletePost(c *fiber.Ctx) error {
	postID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post id"})
	}

	_, err = h.Service.FindById(c.Context(), postID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("Post with id %d not found", postID)})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return h.Service.DeletePost(c.Context(), postID)
}
