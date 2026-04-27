package order_monitor_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
)

type assign interface {
	AssignCourier(ctx context.Context, orderId string) (*model.AssignCourier, error)
}
