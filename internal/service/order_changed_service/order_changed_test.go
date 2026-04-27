package order_changed_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/gateway/order"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/service/order_changed_service/mocks"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestHandleStatusChanged_Success_CallsDo(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := mocks.NewMockOrderStatusFactory(ctrl)
	gw := mocks.NewMockorderGateway(ctrl)
	st := mocks.NewMockOrderStatus(ctrl)

	svc := NewOrderChangedService(f, gw)

	req := model.ChangedStatus{OrderID: "o1", Status: "created"}

	gw.EXPECT().
		GetOrder(gomock.Any(), req.OrderID).
		Return(&order.OrderDto{OrderID: req.OrderID, Status: req.Status}, nil)

	f.EXPECT().Get(req.Status).Return(st)

	st.EXPECT().
		Do(gomock.Any(), req.OrderID).
		Return(nil)

	err := svc.HandleStatusChanged(context.Background(), req)
	require.NoError(t, err)
}

func TestHandleStatusChanged_UnknownStatus(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := mocks.NewMockOrderStatusFactory(ctrl)
	gw := mocks.NewMockorderGateway(ctrl)

	svc := NewOrderChangedService(f, gw)

	req := model.ChangedStatus{OrderID: "o1", Status: "unknown"}

	gw.EXPECT().
		GetOrder(gomock.Any(), req.OrderID).
		Return(&order.OrderDto{OrderID: req.OrderID, Status: req.Status}, nil)

	f.EXPECT().Get(req.Status).Return(nil)

	err := svc.HandleStatusChanged(context.Background(), req)
	require.NoError(t, err)
}

func TestHandleStatusChanged_GatewayError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := mocks.NewMockOrderStatusFactory(ctrl)
	gw := mocks.NewMockorderGateway(ctrl)

	svc := NewOrderChangedService(f, gw)

	req := model.ChangedStatus{OrderID: "o1", Status: "created"}
	Err := errors.New("gw error")

	gw.EXPECT().
		GetOrder(gomock.Any(), req.OrderID).
		Return(nil, Err)

	err := svc.HandleStatusChanged(context.Background(), req)
	require.ErrorIs(t, err, Err)
}

func TestHandleStatusChanged_StatusMismatch(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := mocks.NewMockOrderStatusFactory(ctrl)
	gw := mocks.NewMockorderGateway(ctrl)

	svc := NewOrderChangedService(f, gw)

	req := model.ChangedStatus{OrderID: "o1", Status: "created"}

	gw.EXPECT().
		GetOrder(gomock.Any(), req.OrderID).
		Return(&order.OrderDto{OrderID: req.OrderID, Status: "cancelled"}, nil)

	err := svc.HandleStatusChanged(context.Background(), req)
	require.ErrorIs(t, err, ErrMismatchStatus)
}

func TestHandleStatusChanged_DoReturnsError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := mocks.NewMockOrderStatusFactory(ctrl)
	gw := mocks.NewMockorderGateway(ctrl)
	st := mocks.NewMockOrderStatus(ctrl)

	svc := NewOrderChangedService(f, gw)

	req := model.ChangedStatus{OrderID: "o1", Status: "completed"}
	Err := errors.New("do error")

	gw.EXPECT().
		GetOrder(gomock.Any(), req.OrderID).
		Return(&order.OrderDto{OrderID: req.OrderID, Status: req.Status}, nil)

	f.EXPECT().Get(req.Status).Return(st)

	st.EXPECT().
		Do(gomock.Any(), req.OrderID).
		Return(Err)

	err := svc.HandleStatusChanged(context.Background(), req)
	require.ErrorIs(t, err, Err)
}
