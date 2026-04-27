package order_changed_service

import "errors"

var (
	ErrMismatchStatus = errors.New("request status does not match the current status")
)
