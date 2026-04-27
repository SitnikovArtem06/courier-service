package order_changed_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/service/order_status_factory"
)

type OrderChangedService struct {
	factory order_status_factory.OrderStatusFactory
	gateway orderGateway
}

func NewOrderChangedService(f order_status_factory.OrderStatusFactory, gateway orderGateway) *OrderChangedService {
	return &OrderChangedService{factory: f, gateway: gateway}
}

func (s *OrderChangedService) HandleStatusChanged(ctx context.Context, req model.ChangedStatus) error {

	statusGateway, err := s.gateway.GetOrder(ctx, req.OrderID)
	if err != nil {
		return err
	}

	if req.Status != statusGateway.Status {
		return ErrMismatchStatus
	}

	status := s.factory.Get(req.Status)

	if status == nil {
		return nil
	}

	return status.Do(ctx, req.OrderID)
}
