package courier_service

import "errors"

var (
	ErrDuplicatePhone = errors.New("courier with this phone number already exist")

	ErrNotFound = errors.New("not found")

	ErrInvalidTransport = errors.New("invalid transport")

	ErrInvalidPhoneNumber = errors.New("invalid phone number")

	ErrInvalidStatus = errors.New("invalid status")
)
