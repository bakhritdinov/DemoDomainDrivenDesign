package repository

import (
	"DDD/src/domain"
	"context"
	"errors"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) domain.PostCommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) FindById(ctx context.Context, id int) (*domain.PostComment, error) {
	var comment domain.PostComment
	err := r.db.WithContext(ctx).
		First(&comment, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &comment, err
}

func (r *CommentRepository) FindByPostId(ctx context.Context, postID int) ([]domain.PostComment, error) {
	var comments []domain.PostComment
	err := r.db.WithContext(ctx).
		Where("post_id = ?", postID).
		Find(&comments).Error

	domainComments := make([]domain.PostComment, len(comments))

	return domainComments, err
}

func (r *CommentRepository) Paginate(ctx context.Context, postId int, page int, perPage int) ([]domain.PostComment, int64, error) {
	var comments []domain.PostComment
	var total int64

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.PostComment{}).Where("post_id = ?", postId).Count(&total).Error; err != nil {
			return err
		}

		offset := (page - 1) * perPage
		return tx.
			Where("post_id = ?", postId).
			Order("id DESC").
			Limit(perPage).
			Offset(offset).
			Find(&comments).Error
	})

	return comments, total, err
}

func (r *CommentRepository) Create(ctx context.Context, comment *domain.PostComment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.db.WithContext(ctx).Create(&comment).Error; err != nil {
			return err
		}

		return nil
	})
}
