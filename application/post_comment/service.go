package applicationPostComment

import (
	"DDD/domain"
	"context"
)

type PostCommentService struct {
	PostRepo        domain.PostRepository
	PostCommentRepo domain.PostCommentRepository
}

type PaginatedComments struct {
	Comments   []domain.PostComment `json:"comments"`
	Page       int                  `json:"page"`
	PerPage    int                  `json:"per_page"`
	TotalCount int64                `json:"total_count"`
}

func (s *PostCommentService) FindById(ctx context.Context, commentId int) (*domain.PostComment, error) {
	comment, err := s.PostCommentRepo.FindById(ctx, commentId)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *PostCommentService) FindPaginatedComments(ctx context.Context, postId int, page int, perPage int) (*PaginatedComments, error) {
	comments, total, err := s.PostCommentRepo.Paginate(ctx, postId, page, perPage)
	if err != nil {
		return nil, err
	}

	return &PaginatedComments{
		Comments:   comments,
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
	}, nil
}

func (s *PostCommentService) CreatePostComment(ctx context.Context, comment domain.PostComment) (*domain.PostComment, error) {
	err := s.PostCommentRepo.Create(ctx, &comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}
