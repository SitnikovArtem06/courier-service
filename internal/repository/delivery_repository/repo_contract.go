package delivery_repository

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"time"
)

type DeliveryRepository interface {
	Create(ctx context.Context, orderId string, courierId int64, deadline time.Time) error

	Delete(ctx context.Context, orderId string) (int64, error)

	GetByOrderId(ctx context.Context, orderID string) (*model.DeliveryDB, error)

	GetExpiredOrders(ctx context.Context) ([]int64, error)

	GetCourierWithMinimumOrder(ctx context.Context) (id int64, err error)
}
