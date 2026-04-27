package order_monitor_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/service/assign_service"
	"errors"
	"time"
)

type OrderMonitorService struct {
	gateway  gateway
	assign   assign
	interval time.Duration
	cursor   time.Time
}

func NewOrderMonitorService(gateway gateway, assign assign, interval time.Duration) *OrderMonitorService {
	return &OrderMonitorService{
		gateway:  gateway,
		assign:   assign,
		interval: interval,
		cursor:   time.Now().UTC().Add(-interval),
	}
}

func (s *OrderMonitorService) HandleTick(ctx context.Context) error {
	orders, err := s.gateway.GetNewOrders(ctx, s.cursor)
	if err != nil {
		return err
	}

	for _, id := range orders.OrdersId {
		_, err := s.assign.AssignCourier(ctx, id)
		if err != nil {
			if errors.Is(err, assign_service.ErrNotAvailableCourier) {
				continue
			}
			if errors.Is(err, assign_service.ErrOrderAlreadyAssign) {
				continue
			}
			return err
		}
	}

	var maxCreatedAt = s.cursor
	for _, t := range orders.CreatedAt {
		if t.After(maxCreatedAt) {
			maxCreatedAt = t
		}
	}

	s.cursor = maxCreatedAt
	return nil

}

func (s *OrderMonitorService) Monitor(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.HandleTick(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
