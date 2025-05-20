package value_object

import (
	"fmt"
)

type Text string

func NewText(text string) (Text, error) {
	if err := isValidText(text); err != nil {
		return "", err
	}
	return Text(text), nil
}

func (e Text) String() string {
	return string(e)
}

func isValidText(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("text is required")
	}

	if len(text) < 3 {
		return fmt.Errorf("text is too short")
	}

	if len(text) > 100 {
		return fmt.Errorf("text is too long")
	}
	return nil
}
