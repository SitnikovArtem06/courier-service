package courier_repository

import "errors"

var (
	ErrNotFoundRepo       = errors.New("courier not found")
	ErrDuplicatePhoneRepo = errors.New("duplicate phone")
)
