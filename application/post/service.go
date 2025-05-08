package applicationPost

import (
	"DDD/domain"
	"context"
)

type PostService struct {
	PostRepo domain.PostRepository
}

type PaginatedPosts struct {
	Posts      []domain.Post `json:"posts"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalCount int64         `json:"total_count"`
}

func (s *PostService) FindById(ctx context.Context, postID int) (*domain.Post, error) {
	post, err := s.PostRepo.FindById(ctx, postID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) FindPaginatedPosts(ctx context.Context, page, perPage int) (*PaginatedPosts, error) {
	posts, total, err := s.PostRepo.Paginate(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	return &PaginatedPosts{
		Posts:      posts,
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
	}, nil
}

func (s *PostService) CreatePost(ctx context.Context, post domain.Post) (*domain.Post, error) {
	err := s.PostRepo.Create(ctx, &post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *PostService) UpdatePost(ctx context.Context, post domain.Post) (*domain.Post, error) {
	err := s.PostRepo.Update(ctx, &post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *PostService) DeletePost(ctx context.Context, postID int) error {
	err := s.PostRepo.Delete(ctx, postID)
	if err != nil {
		return err
	}

	return nil
}
