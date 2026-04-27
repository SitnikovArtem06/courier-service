package order_monitor_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/service/assign_service"
	"course-go-avito-SitnikovArtem06/internal/service/order_monitor_service/mocks"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestHandleTick_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gw := mocks.NewMockgateway(ctrl)
	asg := mocks.NewMockassign(ctrl)

	s := NewOrderMonitorService(gw, asg, time.Second)

	startCursor := time.Date(2025, 12, 15, 12, 0, 0, 0, time.UTC)
	s.cursor = startCursor

	t1 := startCursor.Add(10 * time.Second)
	t2 := startCursor.Add(30 * time.Second)
	t3 := startCursor.Add(20 * time.Second)

	resp := &model.OrdersResponse{
		OrdersId:  []string{"o1", "o2", "o3"},
		CreatedAt: []time.Time{t1, t2, t3},
	}

	gw.EXPECT().
		GetNewOrders(gomock.Any(), startCursor).
		Return(resp, nil)

	asg.EXPECT().AssignCourier(gomock.Any(), "o1").Return(nil, nil)
	asg.EXPECT().AssignCourier(gomock.Any(), "o2").Return(nil, nil)
	asg.EXPECT().AssignCourier(gomock.Any(), "o3").Return(nil, nil)

	err := s.HandleTick(context.Background())
	require.NoError(t, err)
	require.Equal(t, t2, s.cursor)
}

func TestHandleTick_Success_EmptyOrders(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gw := mocks.NewMockgateway(ctrl)
	asg := mocks.NewMockassign(ctrl)

	s := NewOrderMonitorService(gw, asg, time.Second)

	startCursor := time.Date(2025, 12, 15, 12, 0, 0, 0, time.UTC)
	s.cursor = startCursor

	resp := &model.OrdersResponse{
		OrdersId:  nil,
		CreatedAt: nil,
	}

	gw.EXPECT().
		GetNewOrders(gomock.Any(), startCursor).
		Return(resp, nil)

	err := s.HandleTick(context.Background())
	require.NoError(t, err)
	require.Equal(t, startCursor, s.cursor)
}

func TestHandleTick_Success_IgnoresNotAvailableCourier(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gw := mocks.NewMockgateway(ctrl)
	asg := mocks.NewMockassign(ctrl)

	s := NewOrderMonitorService(gw, asg, time.Second)

	startCursor := time.Date(2025, 12, 15, 12, 0, 0, 0, time.UTC)
	s.cursor = startCursor

	t1 := startCursor.Add(10 * time.Second)
	resp := &model.OrdersResponse{
		OrdersId:  []string{"o1", "o2"},
		CreatedAt: []time.Time{t1},
	}

	gw.EXPECT().
		GetNewOrders(gomock.Any(), startCursor).
		Return(resp, nil)

	asg.EXPECT().AssignCourier(gomock.Any(), "o1").Return(nil, assign_service.ErrNotAvailableCourier)
	asg.EXPECT().AssignCourier(gomock.Any(), "o2").Return(nil, nil)

	err := s.HandleTick(context.Background())
	require.NoError(t, err)
	require.Equal(t, t1, s.cursor)
}

func TestHandleTick_Success_IgnoresAlreadyAssigned(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gw := mocks.NewMockgateway(ctrl)
	asg := mocks.NewMockassign(ctrl)

	s := NewOrderMonitorService(gw, asg, time.Second)

	startCursor := time.Date(2025, 12, 15, 12, 0, 0, 0, time.UTC)
	s.cursor = startCursor

	t1 := startCursor.Add(10 * time.Second)
	resp := &model.OrdersResponse{
		OrdersId:  []string{"o1", "o2"},
		CreatedAt: []time.Time{t1},
	}

	gw.EXPECT().
		GetNewOrders(gomock.Any(), startCursor).
		Return(resp, nil)

	asg.EXPECT().AssignCourier(gomock.Any(), "o1").Return(nil, assign_service.ErrOrderAlreadyAssign)
	asg.EXPECT().AssignCourier(gomock.Any(), "o2").Return(nil, nil)

	err := s.HandleTick(context.Background())
	require.NoError(t, err)
	require.Equal(t, t1, s.cursor)
}

func TestHandleTick_GatewayError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gw := mocks.NewMockgateway(ctrl)
	asg := mocks.NewMockassign(ctrl)

	s := NewOrderMonitorService(gw, asg, time.Second)

	startCursor := time.Date(2025, 12, 15, 12, 0, 0, 0, time.UTC)
	s.cursor = startCursor

	gwErr := errors.New("gateway error")

	gw.EXPECT().
		GetNewOrders(gomock.Any(), startCursor).
		Return(nil, gwErr)

	err := s.HandleTick(context.Background())
	require.ErrorIs(t, err, gwErr)
	require.Equal(t, startCursor, s.cursor)
}

func TestHandleTick_Error_AssignError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gw := mocks.NewMockgateway(ctrl)
	asg := mocks.NewMockassign(ctrl)

	s := NewOrderMonitorService(gw, asg, time.Second)

	startCursor := time.Date(2025, 12, 15, 12, 0, 0, 0, time.UTC)
	s.cursor = startCursor

	t1 := startCursor.Add(10 * time.Second)
	resp := &model.OrdersResponse{
		OrdersId:  []string{"o1", "o2"},
		CreatedAt: []time.Time{t1},
	}

	gw.EXPECT().
		GetNewOrders(gomock.Any(), startCursor).
		Return(resp, nil)

	asgErr := errors.New("assign error")

	asg.EXPECT().AssignCourier(gomock.Any(), "o1").Return(nil, asgErr)

	err := s.HandleTick(context.Background())
	require.ErrorIs(t, err, asgErr)
	require.Equal(t, startCursor, s.cursor)
}
