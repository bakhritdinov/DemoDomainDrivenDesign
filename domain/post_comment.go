package domain

import (
	"context"
	"time"
)

type PostComment struct {
	Id        uint       `json:"id"`
	PostId    uint       `json:"postId"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type PostCommentRepository interface {
	FindById(ctx context.Context, id int) (*PostComment, error)
	FindByPostId(ctx context.Context, postID int) ([]PostComment, error)
	Paginate(ctx context.Context, postId int, page int, perPage int) ([]PostComment, int64, error)
	Create(ctx context.Context, comment *PostComment) error
}
