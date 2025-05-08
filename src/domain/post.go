package domain

import (
	"context"
	"time"
)

type Post struct {
	Id        uint
	Title     string
	Content   string
	Comments  []PostComment `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type PostRepository interface {
	FindById(ctx context.Context, id int) (*Post, error)
	Paginate(ctx context.Context, page int, perPage int) ([]Post, int64, error)
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int) error
}
