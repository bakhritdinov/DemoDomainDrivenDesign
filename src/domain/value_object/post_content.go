package value_object

import "fmt"

type PostContent struct {
	Value string `json:"value"`
}

func NewPostContent(value string) (PostContent, error) {
	if err := ValidatePostContent(value); err != nil {
		return PostContent{}, err
	}
	return PostContent{Value: value}, nil
}

func (p PostTitle) Equals(other PostTitle) bool {
	return p.Value == other.Value
}

func ValidatePostContent(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("post content is required")
	}

	if len(value) > 500 {
		return fmt.Errorf("post content is too long")
	}

	return nil
}
