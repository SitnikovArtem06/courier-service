package delivery_repository

import "errors"

var (
	ErrNotFound    error = errors.New("delivery not found")
	ErrNoneCourier error = errors.New("no one courier can be assigned")
)
