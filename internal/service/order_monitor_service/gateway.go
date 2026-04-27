package order_monitor_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"time"
)

type gateway interface {
	GetNewOrders(ctx context.Context, from time.Time) (*model.OrdersResponse, error)
}
