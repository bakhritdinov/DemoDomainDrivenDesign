package value_object

import (
	"errors"
)

type Title string

func NewTitle(title string) (Title, error) {
	if err := isValidTitle(title); err != nil {
		return "", errors.New(err.Error())
	}

	return Title(title), nil
}

func (e Title) String() string {
	return string(e)
}

func isValidTitle(title string) error {
	if len(title) == 0 {
		return errors.New("title is required")
	}

	if len(title) < 3 {
		return errors.New("title is too short")
	}

	if len(title) > 100 {
		return errors.New("title is too long")
	}
	return nil
}
