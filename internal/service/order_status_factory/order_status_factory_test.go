package order_status_factory

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/service/order_status_factory/mocks"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestOrderStatusFactory_Get_UnknownNil(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := mocks.NewMockassign(ctrl)
	f := NewOrderStatusFactory(a)

	require.Nil(t, f.Get("unknown"))
}

func TestCreated_Do_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := mocks.NewMockassign(ctrl)
	f := NewOrderStatusFactory(a)

	orderId := "o1"

	a.EXPECT().
		AssignCourier(gomock.Any(), orderId).
		Return(&model.AssignCourier{OrderId: orderId}, nil)

	err := f.Get("created").Do(context.Background(), orderId)
	require.NoError(t, err)
}

func TestCreated_Do_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := mocks.NewMockassign(ctrl)
	f := NewOrderStatusFactory(a)

	orderId := "o1"
	Err := errors.New("assign err")

	a.EXPECT().
		AssignCourier(gomock.Any(), orderId).
		Return(nil, Err)

	err := f.Get("created").Do(context.Background(), orderId)
	require.ErrorIs(t, err, Err)
}

func TestCancelled_Do_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := mocks.NewMockassign(ctrl)
	f := NewOrderStatusFactory(a)

	orderId := "o1"

	a.EXPECT().
		UnassignCourier(gomock.Any(), orderId).
		Return(&model.UnassignCourier{OrderId: orderId}, nil)

	err := f.Get("cancelled").Do(context.Background(), orderId)
	require.NoError(t, err)
}

func TestCancelled_Do_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := mocks.NewMockassign(ctrl)
	f := NewOrderStatusFactory(a)

	orderId := "o1"
	Err := errors.New("unassign err")

	a.EXPECT().
		UnassignCourier(gomock.Any(), orderId).
		Return(nil, Err)

	err := f.Get("cancelled").Do(context.Background(), orderId)
	require.ErrorIs(t, err, Err)
}

func TestCompleted_Do_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := mocks.NewMockassign(ctrl)
	f := NewOrderStatusFactory(a)

	orderId := "o1"

	a.EXPECT().
		CompleteCourier(gomock.Any(), orderId).
		Return(nil)

	err := f.Get("completed").Do(context.Background(), orderId)
	require.NoError(t, err)
}

func TestCompleted_Do_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := mocks.NewMockassign(ctrl)
	f := NewOrderStatusFactory(a)

	orderId := "o1"
	Err := errors.New("complete err")

	a.EXPECT().
		CompleteCourier(gomock.Any(), orderId).
		Return(Err)

	err := f.Get("completed").Do(context.Background(), orderId)
	require.ErrorIs(t, err, Err)
}
