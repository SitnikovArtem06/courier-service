package courier_handler

import "errors"

var (
	ErrInvalidId = errors.New("invalid ID")
	ErrEmptyName = errors.New("name is empty")
)
