package domain

import (
	"DDD/src/domain/value_object"
	"context"
	"time"
)

type PostComment struct {
	Id        uint              `gorm:"primarykey" json:"id"`
	PostId    uint              `gorm:"index;not null" json:"postId"`
	Text      value_object.Text `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	DeletedAt *time.Time        `gorm:"index" json:"deletedAt"`
}

type PostCommentRepository interface {
	FindById(ctx context.Context, id int) (*PostComment, error)
	FindByPostId(ctx context.Context, postID int) ([]PostComment, error)
	Paginate(ctx context.Context, postId int, page int, perPage int) ([]PostComment, int64, error)
	Create(ctx context.Context, comment *PostComment) error
}
