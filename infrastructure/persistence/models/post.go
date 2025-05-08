package models

import (
	"DDD/domain"
	"gorm.io/gorm"
)

type PostTable struct {
	gorm.Model
	Title    string             `gorm:"size:255;not null" validate:"required,min=3,max=100"`
	Content  string             `gorm:"type:text" validate:"required,max=500"`
	Comments []PostCommentTable `gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE" validate:"dive"`
}

func (PostTable) TableName() string {
	return "posts"
}

func ToDomainPost(p PostTable) domain.Post {
	comments := make([]domain.PostComment, len(p.Comments))

	for i, comment := range p.Comments {
		domainComment := domain.PostComment{
			Id:        comment.ID,
			Text:      comment.Text,
			PostId:    p.ID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}

		if comment.DeletedAt.Valid {
			domainComment.DeletedAt = &comment.DeletedAt.Time
		} else {
			domainComment.DeletedAt = nil
		}

		comments[i] = domainComment
	}

	domainPost := domain.Post{
		Id:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		Comments:  comments,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	if p.DeletedAt.Valid {
		domainPost.DeletedAt = &p.DeletedAt.Time
	} else {
		domainPost.DeletedAt = nil
	}

	return domainPost
}

func FromPostDomain(p domain.Post) PostTable {
	comments := make([]PostCommentTable, len(p.Comments))

	for i, comment := range p.Comments {
		commentTable := PostCommentTable{
			Model: gorm.Model{
				ID:        comment.Id,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			Text:   comment.Text,
			PostId: p.Id,
		}

		if comment.DeletedAt != nil {
			commentTable.DeletedAt = gorm.DeletedAt{Time: *comment.DeletedAt, Valid: true}
		} else {
			commentTable.DeletedAt = gorm.DeletedAt{}
		}

		comments[i] = commentTable
	}

	postTable := PostTable{
		Model: gorm.Model{
			ID:        p.Id,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
		Title:    p.Title,
		Content:  p.Content,
		Comments: comments,
	}

	if p.DeletedAt != nil {
		postTable.Model.DeletedAt = gorm.DeletedAt{Time: *p.DeletedAt, Valid: true}
	} else {
		postTable.Model.DeletedAt = gorm.DeletedAt{}
	}

	return postTable
}
