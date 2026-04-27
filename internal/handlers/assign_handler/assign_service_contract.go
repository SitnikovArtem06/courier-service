package assign_handler

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
)

type assignService interface {
	AssignCourier(ctx context.Context, orderId string) (*model.AssignCourier, error)

	UnassignCourier(ctx context.Context, orderId string) (*model.UnassignCourier, error)
}
