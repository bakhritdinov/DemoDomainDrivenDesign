package value_object

import (
	"errors"
)

type Content string

func NewContent(content string) (Content, error) {
	if err := isValidContent(content); err != nil {
		return "", errors.New(err.Error())
	}

	return Content(content), nil
}

func (e Content) String() string {
	return string(e)
}

func isValidContent(content string) error {
	if len(content) == 0 {
		return errors.New("content is required")
	}

	if len(content) > 500 {
		return errors.New("content is too long")
	}

	return nil
}
