package delivery_monitor_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/repository/courier_repository"
	"course-go-avito-SitnikovArtem06/internal/repository/delivery_repository"
	"time"
)

type DeliveryMonitorService struct {
	dRepo    delivery_repository.DeliveryRepository
	cRepo    courier_repository.CourierRepository
	interval time.Duration
}

func NewDeliveryMonitorService(dRepo delivery_repository.DeliveryRepository, cRepo courier_repository.CourierRepository, interval time.Duration) *DeliveryMonitorService {
	return &DeliveryMonitorService{dRepo: dRepo, cRepo: cRepo, interval: interval}
}

func (s *DeliveryMonitorService) handleTick(ctx context.Context) error {
	ids, err := s.dRepo.GetExpiredOrders(ctx)
	if err != nil {
		return err
	}

	if len(ids) > 0 {
		if err := s.cRepo.UpdateAllExpiredCourier(ctx, ids); err != nil {
			return err
		}
	}

	return nil
}

func (s *DeliveryMonitorService) MonitorDeadline(ctx context.Context) error {

	ticker := time.NewTicker(s.interval)

	defer ticker.Stop()

	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:

			if err := s.handleTick(ctx); err != nil {
				return err
			}

		}
	}

}
