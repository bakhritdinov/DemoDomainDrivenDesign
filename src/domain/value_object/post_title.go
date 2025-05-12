package value_object

import (
	"fmt"
)

type PostTitle struct {
	Value string `json:"value"`
}

func NewPostTitle(value string) (PostTitle, error) {
	if err := ValidatePostTitle(value); err != nil {
		return PostTitle{}, err
	}
	return PostTitle{Value: value}, nil
}

func ValidatePostTitle(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("post title is required")
	}

	if len(value) < 3 {
		return fmt.Errorf("post title is too short")
	}

	if len(value) > 100 {
		return fmt.Errorf("post title is too long")
	}
	return nil
}
