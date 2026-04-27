package assign_service

import "errors"

var (
	ErrNotAvailableCourier error = errors.New("no available couriers now")

	ErrOrderAlreadyAssign error = errors.New("order already assign")

	ErrNotAssignedCourier error = errors.New("no one courier associated with this order")

	ErrNotFoundOrder = errors.New("not found order")
)
