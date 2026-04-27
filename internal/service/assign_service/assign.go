package assign_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/repository/courier_repository"
	"course-go-avito-SitnikovArtem06/internal/repository/delivery_repository"
	"course-go-avito-SitnikovArtem06/internal/service/transport_factory"
	"course-go-avito-SitnikovArtem06/internal/tx"
	"errors"
)

type AssignService struct {
	txManager        tx.TransactionManager
	deliveryRepo     delivery_repository.DeliveryRepository
	courierRepo      courier_repository.CourierRepository
	TransportFactory transport_factory.TransportFactory
}

func NewAssignService(txManager tx.TransactionManager, dRepo delivery_repository.DeliveryRepository, cRepo courier_repository.CourierRepository, transportF transport_factory.TransportFactory) *AssignService {
	return &AssignService{
		txManager:        txManager,
		deliveryRepo:     dRepo,
		courierRepo:      cRepo,
		TransportFactory: transportF,
	}
}

func (s *AssignService) AssignCourier(ctx context.Context, orderId string) (*model.AssignCourier, error) {

	var result *model.AssignCourier

	err := s.txManager.Begin(ctx, true, func(ctx context.Context) error {

		if _, err := s.deliveryRepo.GetByOrderId(ctx, orderId); err == nil {
			return ErrOrderAlreadyAssign
		} else if !errors.Is(err, delivery_repository.ErrNotFound) {
			return err
		}

		courierId, err := s.deliveryRepo.GetCourierWithMinimumOrder(ctx)
		if err != nil {
			if errors.Is(err, delivery_repository.ErrNoneCourier) {
				return ErrNotAvailableCourier
			}
			return err
		}

		courier, err := s.courierRepo.Get(ctx, courierId)
		if err != nil {
			return err
		}

		tr := s.TransportFactory.Get(courier.Transport)
		deadline := tr.Deadline()

		if err = s.deliveryRepo.Create(ctx, orderId, courier.Id, deadline); err != nil {
			return err
		}

		status := model.CourierStatusBusy

		req := &model.UpdateCourierRequest{Id: &courier.Id, Status: &status}

		if err = s.courierRepo.Update(ctx, req); err != nil {
			return err
		}

		result = &model.AssignCourier{
			CourierId: courier.Id,
			OrderId:   orderId,
			Transport: courier.Transport,
			Deadline:  deadline,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil

}

func (s *AssignService) UnassignCourier(ctx context.Context, orderId string) (*model.UnassignCourier, error) {

	var unassign *model.UnassignCourier

	err := s.txManager.Begin(ctx, true, func(ctx context.Context) error {

		courierId, err := s.deliveryRepo.Delete(ctx, orderId)
		if err != nil {
			if errors.Is(err, delivery_repository.ErrNotFound) {
				return ErrNotAssignedCourier
			}
			return err
		}

		status := model.CourierStatusAvailable

		if err = s.courierRepo.Update(ctx, &model.UpdateCourierRequest{Id: &courierId, Status: &status}); err != nil {
			return err
		}

		unassign = &model.UnassignCourier{
			CourierId: courierId,
			OrderId:   orderId,
			Status:    model.Unassigned,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return unassign, nil

}

func (s *AssignService) CompleteCourier(ctx context.Context, orderId string) error {
	err := s.txManager.Begin(ctx, true, func(ctx context.Context) error {
		delivery, err := s.deliveryRepo.GetByOrderId(ctx, orderId)
		if err != nil {
			if errors.Is(err, delivery_repository.ErrNotFound) {
				return ErrNotFoundOrder
			}
			return err
		}

		status := model.CourierStatusAvailable
		reqUpdate := &model.UpdateCourierRequest{
			Id:     &delivery.CourierId,
			Status: &status,
		}
		err = s.courierRepo.Update(ctx, reqUpdate)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
