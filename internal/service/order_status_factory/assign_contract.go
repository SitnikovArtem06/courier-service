package order_status_factory

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
)

type assign interface {
	AssignCourier(ctx context.Context, orderId string) (*model.AssignCourier, error)
	UnassignCourier(ctx context.Context, orderId string) (*model.UnassignCourier, error)

	CompleteCourier(ctx context.Context, orderId string) error
}
