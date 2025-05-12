package models

import (
	"DDD/src/domain"
	"DDD/src/domain/value_object"
	"gorm.io/gorm"
)

type PostCommentTable struct {
	gorm.Model
	PostId uint   `gorm:"index;not null"`
	Text   string `gorm:"type:text;not null"`
}

func (PostCommentTable) TableName() string {
	return "post_comments"
}

func ToDomainPostComment(p PostCommentTable) domain.PostComment {

	postCommentText, err := value_object.NewPostCommentText(p.Text)
	if err != nil {
		panic(err)
	}

	domainComment := domain.PostComment{
		Id:        p.ID,
		Text:      postCommentText,
		PostId:    p.PostId,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	if p.DeletedAt.Valid {
		domainComment.DeletedAt = &p.DeletedAt.Time
	} else {
		domainComment.DeletedAt = nil
	}

	return domainComment
}

func FromPostCommentDomain(p domain.PostComment) PostCommentTable {
	commentTable := PostCommentTable{
		Model: gorm.Model{
			ID:        p.Id,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
		PostId: p.PostId,
		Text:   p.Text.Value,
	}

	if p.DeletedAt != nil {
		commentTable.Model.DeletedAt = gorm.DeletedAt{Time: *p.DeletedAt, Valid: true}
	} else {
		commentTable.Model.DeletedAt = gorm.DeletedAt{}
	}

	return commentTable
}
