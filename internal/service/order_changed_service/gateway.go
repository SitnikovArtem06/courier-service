package order_changed_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/gateway/order"
)

type orderGateway interface {
	GetOrder(ctx context.Context, orderID string) (*order.OrderDto, error)
}
