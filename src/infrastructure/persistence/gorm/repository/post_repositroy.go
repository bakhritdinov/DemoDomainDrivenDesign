package repository

import (
	"DDD/src/domain"
	"context"
	"errors"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) domain.PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) FindById(ctx context.Context, id int) (*domain.Post, error) {
	var post domain.Post
	err := r.db.WithContext(ctx).
		First(&post, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &post, err
}

func (r *PostRepository) Paginate(ctx context.Context, page int, perPage int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Post{}).Count(&total).Error; err != nil {
			return err
		}

		offset := (page - 1) * perPage
		return tx.
			Order("id DESC").
			Limit(perPage).
			Offset(offset).
			Find(&posts).Error
	})

	return posts, total, err
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.Post
		if err := tx.Where("title = ?", post.Title).First(&existing).Error; err == nil {
			return gorm.ErrDuplicatedKey
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := r.db.WithContext(ctx).Create(&post).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.Post
		if err := tx.Where("title = ? AND id != ?", post.Title, post.Id).First(&existing).Error; err == nil {
			return gorm.ErrDuplicatedKey
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := r.db.WithContext(ctx).Updates(&post).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PostRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var post domain.Post
		if err := tx.First(&post, id).Error; err != nil {
			return err
		}

		return tx.Delete(&post).Error
	})
}
