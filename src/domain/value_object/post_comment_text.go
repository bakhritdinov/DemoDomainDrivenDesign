package value_object

import (
	"fmt"
)

type PostCommentText struct {
	Value string `json:"value"`
}

func NewPostCommentText(value string) (PostCommentText, error) {
	if err := ValidatePostCommentText(value); err != nil {
		return PostCommentText{}, err
	}
	return PostCommentText{Value: value}, nil
}

func ValidatePostCommentText(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("post comment text is required")
	}

	if len(value) < 3 {
		return fmt.Errorf("post comment text is too short")
	}

	if len(value) > 100 {
		return fmt.Errorf("post comment text is too long")
	}
	return nil
}
