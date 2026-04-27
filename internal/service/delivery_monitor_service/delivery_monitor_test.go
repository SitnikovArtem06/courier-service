package delivery_monitor_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/service/mocks"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestHandleTick_NoExpiredOrders(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	cRepo := mocks.NewMockCourierRepository(ctrl)

	service := &DeliveryMonitorService{
		dRepo:    dRepo,
		cRepo:    cRepo,
		interval: 0,
	}

	dRepo.EXPECT().
		GetExpiredOrders(gomock.Any()).
		Return([]int64{}, nil).
		Times(1)

	cRepo.EXPECT().
		UpdateAllExpiredCourier(gomock.Any(), gomock.Any()).
		Times(0)

	err := service.handleTick(context.Background())
	require.NoError(t, err)
}

func TestHandleTick_WithExpiredOrders(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	cRepo := mocks.NewMockCourierRepository(ctrl)

	service := &DeliveryMonitorService{
		dRepo:    dRepo,
		cRepo:    cRepo,
		interval: 0,
	}

	ids := []int64{1, 2, 3}

	dRepo.EXPECT().
		GetExpiredOrders(gomock.Any()).
		Return(ids, nil).
		Times(1)

	cRepo.EXPECT().
		UpdateAllExpiredCourier(gomock.Any(), ids).
		Return(nil).
		Times(1)

	err := service.handleTick(context.Background())
	require.NoError(t, err)
}

func TestHandleTick_DBError_GetExpiredOrders(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	cRepo := mocks.NewMockCourierRepository(ctrl)

	service := &DeliveryMonitorService{
		dRepo:    dRepo,
		cRepo:    cRepo,
		interval: 0,
	}
	context.Background()
	dbErr := errors.New("db error")

	dRepo.EXPECT().
		GetExpiredOrders(gomock.Any()).
		Return(nil, dbErr).
		Times(1)

	cRepo.EXPECT().
		UpdateAllExpiredCourier(gomock.Any(), gomock.Any()).
		Times(0)

	err := service.handleTick(context.Background())
	require.ErrorIs(t, err, dbErr)
}

func TestHandleTick_DBError_UpdateAllExpiredCourier(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	cRepo := mocks.NewMockCourierRepository(ctrl)

	service := &DeliveryMonitorService{
		dRepo:    dRepo,
		cRepo:    cRepo,
		interval: 0,
	}

	ids := []int64{1, 2, 3}
	wantErr := errors.New("db error")

	dRepo.EXPECT().
		GetExpiredOrders(gomock.Any()).
		Return(ids, nil).
		Times(1)

	cRepo.EXPECT().
		UpdateAllExpiredCourier(gomock.Any(), ids).
		Return(wantErr).
		Times(1)

	err := service.handleTick(context.Background())
	require.ErrorIs(t, err, wantErr)
}

func TestMonitorDeadline_ContextTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	cRepo := mocks.NewMockCourierRepository(ctrl)

	interval := time.Millisecond

	service := NewDeliveryMonitorService(dRepo, cRepo, interval)

	dRepo.EXPECT().
		GetExpiredOrders(gomock.Any()).
		Return([]int64{}, nil).
		AnyTimes()

	cRepo.EXPECT().
		UpdateAllExpiredCourier(gomock.Any(), gomock.Any()).
		Times(0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- service.MonitorDeadline(ctx)
	}()

	err := <-errCh
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestMonitorDeadline_HandleTickError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	cRepo := mocks.NewMockCourierRepository(ctrl)

	interval := time.Millisecond

	service := NewDeliveryMonitorService(dRepo, cRepo, interval)

	wantErr := errors.New("get expired error")

	dRepo.EXPECT().
		GetExpiredOrders(gomock.Any()).
		Return(nil, wantErr).
		Times(1)

	cRepo.EXPECT().
		UpdateAllExpiredCourier(gomock.Any(), gomock.Any()).
		Times(0)

	ctx := context.Background()

	errCh := make(chan error, 1)
	go func() {
		errCh <- service.MonitorDeadline(ctx)
	}()

	err := <-errCh
	require.ErrorIs(t, err, wantErr)
}
