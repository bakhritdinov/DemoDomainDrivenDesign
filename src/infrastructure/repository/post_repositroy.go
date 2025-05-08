package repository

import (
	"DDD/src/domain"
	"DDD/src/infrastructure/persistence/models"
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
	var postTable models.PostTable
	err := r.db.WithContext(ctx).
		First(&postTable, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	post := models.ToDomainPost(postTable)

	return &post, err
}

func (r *PostRepository) Paginate(ctx context.Context, page int, perPage int) ([]domain.Post, int64, error) {
	var posts []models.PostTable
	var total int64

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.PostTable{}).Count(&total).Error; err != nil {
			return err
		}

		offset := (page - 1) * perPage
		return tx.
			Order("id DESC").
			Limit(perPage).
			Offset(offset).
			Find(&posts).Error
	})

	domainPosts := make([]domain.Post, len(posts))
	for i, p := range posts {
		domainPosts[i] = models.ToDomainPost(p)
	}

	return domainPosts, total, err
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing models.PostTable
		if err := tx.Where("title = ?", post.Title).First(&existing).Error; err == nil {
			return gorm.ErrDuplicatedKey
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		postTable := models.FromPostDomain(*post)

		if err := r.db.WithContext(ctx).Create(&postTable).Error; err != nil {
			return err
		}

		*post = models.ToDomainPost(postTable) // update domain model
		return nil
	})
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing models.PostTable
		if err := tx.Where("title = ?", post.Title).First(&existing).Error; err == nil {
			return gorm.ErrDuplicatedKey
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		postTable := models.FromPostDomain(*post)

		if err := r.db.WithContext(ctx).Updates(&postTable).Error; err != nil {
			return err
		}

		*post = models.ToDomainPost(postTable) // update domain model
		return nil
	})
}

func (r *PostRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var post models.PostTable
		if err := tx.First(&post, id).Error; err != nil {
			return err
		}

		return tx.Delete(&post).Error
	})
}
