package repository

import (
	"DDD/domain"
	"DDD/infrastructure/persistence/models"
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
	var commentTable models.PostCommentTable
	err := r.db.WithContext(ctx).
		First(&commentTable, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	comment := models.ToDomainPostComment(commentTable)

	return &comment, err
}

func (r *CommentRepository) FindByPostId(ctx context.Context, postID int) ([]domain.PostComment, error) {
	var comments []models.PostCommentTable
	err := r.db.WithContext(ctx).
		Where("post_id = ?", postID).
		Find(&comments).Error

	domainComments := make([]domain.PostComment, len(comments))
	for i, c := range comments {
		domainComments[i] = models.ToDomainPostComment(c)
	}

	return domainComments, err
}

func (r *CommentRepository) Paginate(ctx context.Context, postId int, page int, perPage int) ([]domain.PostComment, int64, error) {
	var comments []models.PostCommentTable
	var total int64

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.PostCommentTable{}).Where("post_id = ?", postId).Count(&total).Error; err != nil {
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

	domainComments := make([]domain.PostComment, len(comments))
	for i, c := range comments {
		domainComments[i] = models.ToDomainPostComment(c)
	}

	return domainComments, total, err
}

func (r *CommentRepository) Create(ctx context.Context, comment *domain.PostComment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		commentTable := models.FromPostCommentDomain(*comment)

		if err := r.db.WithContext(ctx).Create(&commentTable).Error; err != nil {
			return err
		}

		*comment = models.ToDomainPostComment(commentTable) // update domain model
		return nil
	})
}
