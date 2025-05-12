package models

import (
	"DDD/src/domain"
	"DDD/src/domain/value_object"
	"gorm.io/gorm"
)

type PostTable struct {
	gorm.Model
	Title    string             `gorm:"size:255;not null"`
	Content  string             `gorm:"type:text"`
	Comments []PostCommentTable `gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE"`
}

func (PostTable) TableName() string {
	return "posts"
}

func ToDomainPost(p PostTable) domain.Post {
	comments := make([]domain.PostComment, len(p.Comments))

	for i, comment := range p.Comments {
		postCommentText, err := value_object.NewPostCommentText(comment.Text)
		if err != nil {
			panic(err)
		}

		domainComment := domain.PostComment{
			Id:        comment.ID,
			Text:      postCommentText,
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

	postTitle, err := value_object.NewPostTitle(p.Title)
	if err != nil {
		panic(err)
	}

	postContent, err := value_object.NewPostContent(p.Content)
	if err != nil {
		panic(err)
	}

	domainPost := domain.Post{
		Id:        p.ID,
		Title:     postTitle,
		Content:   postContent,
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
			Text:   comment.Text.Value,
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
		Title:    p.Title.Value,
		Content:  p.Content.Value,
		Comments: comments,
	}

	if p.DeletedAt != nil {
		postTable.Model.DeletedAt = gorm.DeletedAt{Time: *p.DeletedAt, Valid: true}
	} else {
		postTable.Model.DeletedAt = gorm.DeletedAt{}
	}

	return postTable
}
