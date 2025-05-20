package domain

import (
	"DDD/src/domain/value_object"
	"context"
	"time"
)

type Post struct {
	Id        uint                 `gorm:"primarykey" json:"id"`
	Title     value_object.Title   `gorm:"size:255;not null" json:"title"`
	Content   value_object.Content `gorm:"type:text" json:"content"`
	Comments  []PostComment        `gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	DeletedAt *time.Time           `gorm:"index" json:"deletedAt"`
}

type PostRepository interface {
	FindById(ctx context.Context, id int) (*Post, error)
	Paginate(ctx context.Context, page int, perPage int) ([]Post, int64, error)
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int) error
}
