package domain

import (
	"context"
	"time"
)

type Post struct {
	Id        uint          `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	Comments  []PostComment `json:"comments,omitempty"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	DeletedAt *time.Time    `json:"deletedAt"`
}

type PostRepository interface {
	FindById(ctx context.Context, id int) (*Post, error)
	Paginate(ctx context.Context, page int, perPage int) ([]Post, int64, error)
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int) error
}
