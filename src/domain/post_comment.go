package domain

import (
	"context"
	"time"
)

type PostComment struct {
	Id        uint
	PostId    uint
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type PostCommentRepository interface {
	FindById(ctx context.Context, id int) (*PostComment, error)
	FindByPostId(ctx context.Context, postID int) ([]PostComment, error)
	Paginate(ctx context.Context, postId int, page int, perPage int) ([]PostComment, int64, error)
	Create(ctx context.Context, comment *PostComment) error
}
