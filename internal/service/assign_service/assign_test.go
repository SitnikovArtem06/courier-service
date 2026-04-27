package assign_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/repository/delivery_repository"
	"course-go-avito-SitnikovArtem06/internal/service/mocks"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestAssign_Success(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)

	tx := mocks.NewMockTransactionManager(ctrl)

	dRepo := mocks.NewMockDeliveryRepository(ctrl)

	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderId).
		Return(nil, delivery_repository.ErrNotFound)

	dRepo.EXPECT().GetCourierWithMinimumOrder(gomock.Any()).Return(int64(1), nil)

	courier := &model.CourierDB{
		Id:        1,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.Car,
	}

	cRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(courier, nil)

	deadline := time.Now().Add(15 * time.Minute).UTC()

	tMock := mocks.NewMockTransport(ctrl)

	transportFactory.EXPECT().
		Get(courier.Transport).
		Return(tMock)

	tMock.EXPECT().
		Deadline().
		Return(deadline)

	dRepo.EXPECT().
		Create(gomock.Any(), orderId, int64(1), deadline).
		Return(nil)

	cRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, r *model.UpdateCourierRequest) error {
			require.NotNil(t, r.Id)
			require.NotNil(t, r.Status)
			require.Equal(t, int64(1), *r.Id)
			require.Equal(t, model.CourierStatusBusy, *r.Status)
			return nil
		})

	got, err := service.AssignCourier(context.Background(), orderId)

	require.NoError(t, err)
	require.Equal(t, int64(1), got.CourierId)
	require.Equal(t, orderId, got.OrderId)
	require.Equal(t, model.Car, got.Transport)
	require.Equal(t, deadline, got.Deadline)
}

func TestAssign_AlreadyAssign(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)

	tx := mocks.NewMockTransactionManager(ctrl)

	dRepo := mocks.NewMockDeliveryRepository(ctrl)

	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"

	dRepo.EXPECT().GetByOrderId(gomock.Any(), orderId).Return(nil, nil)

	got, err := service.AssignCourier(context.Background(), orderId)

	require.Nil(t, got)
	require.ErrorIs(t, ErrOrderAlreadyAssign, err)

}

func TestAssign_NobodyAvailable(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)

	tx := mocks.NewMockTransactionManager(ctrl)

	dRepo := mocks.NewMockDeliveryRepository(ctrl)

	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderId).
		Return(nil, delivery_repository.ErrNotFound)

	var couriersExpected []model.CourierDB

	dRepo.EXPECT().GetCourierWithMinimumOrder(gomock.Any()).Return(int64(0), delivery_repository.ErrNoneCourier)

	got, err := service.AssignCourier(context.Background(), orderId)

	require.Nil(t, got)
	require.Equal(t, len(couriersExpected), 0)
	require.ErrorIs(t, ErrNotAvailableCourier, err)

}

func TestAssign_DBError_GetByOrderId(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)
	orderID := "1"
	dbErr := errors.New("db error")

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderID).
		Return(nil, dbErr)

	got, err := service.AssignCourier(context.Background(), orderID)

	require.Nil(t, got)
	require.Equal(t, dbErr, err)
}
func TestAssign_DBError_GetCourierWithMinimumOrder(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)
	orderID := "1"
	dbErr := errors.New("db error")

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderID).
		Return(nil, delivery_repository.ErrNotFound)

	dRepo.EXPECT().
		GetCourierWithMinimumOrder(gomock.Any()).
		Return(int64(0), dbErr)

	got, err := service.AssignCourier(context.Background(), orderID)

	require.Nil(t, got)
	require.Equal(t, dbErr, err)
}

func TestAssign_DBError_GetCourier(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	orderID := "1"
	dbErr := errors.New("db error")

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderID).
		Return(nil, delivery_repository.ErrNotFound)

	dRepo.EXPECT().
		GetCourierWithMinimumOrder(gomock.Any()).
		Return(int64(1), nil)

	cRepo.EXPECT().
		Get(gomock.Any(), int64(1)).
		Return(nil, dbErr)

	got, err := service.AssignCourier(context.Background(), orderID)

	require.Nil(t, got)
	require.Equal(t, dbErr, err)
}

func TestAssign_DBError_CreateDelivery(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)
	orderID := "1"
	dbErr := errors.New("db error")

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderID).
		Return(nil, delivery_repository.ErrNotFound)

	dRepo.EXPECT().
		GetCourierWithMinimumOrder(gomock.Any()).
		Return(int64(1), nil)

	courier := &model.CourierDB{
		Id:        1,
		Transport: model.Car,
	}

	cRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(courier, nil)

	deadline := time.Now().UTC()
	tMock := mocks.NewMockTransport(ctrl)
	transportFactory.EXPECT().
		Get(courier.Transport).
		Return(tMock)
	tMock.EXPECT().
		Deadline().
		Return(deadline)

	dRepo.EXPECT().
		Create(gomock.Any(), orderID, int64(1), deadline).
		Return(dbErr)

	got, err := service.AssignCourier(context.Background(), orderID)

	require.Nil(t, got)
	require.Equal(t, dbErr, err)
}

func TestAssign_DBError_UpdateCourier(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)
	orderID := "1"
	dbErr := errors.New("db error")

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderID).
		Return(nil, delivery_repository.ErrNotFound)

	dRepo.EXPECT().
		GetCourierWithMinimumOrder(gomock.Any()).
		Return(int64(1), nil)

	courier := &model.CourierDB{
		Id:        1,
		Transport: model.Car,
	}

	cRepo.EXPECT().Get(gomock.Any(), int64(1)).Return(courier, nil)

	deadline := time.Now().UTC()
	tMock := mocks.NewMockTransport(ctrl)
	transportFactory.EXPECT().
		Get(model.Car).
		Return(tMock)
	tMock.EXPECT().
		Deadline().
		Return(deadline)

	dRepo.EXPECT().
		Create(gomock.Any(), orderID, int64(1), deadline).
		Return(nil)

	cRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(dbErr)

	got, err := service.AssignCourier(context.Background(), orderID)

	require.Nil(t, got)
	require.Equal(t, dbErr, err)
}

func TestAssign_DBError_Begin(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)
	orderID := "1"
	dbErr := errors.New("tx begin error")

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		Return(dbErr)

	got, err := service.AssignCourier(context.Background(), orderID)

	require.Nil(t, got)
	require.Equal(t, dbErr, err)
}

func TestUnassign_Success(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"

	dModel := &model.DeliveryDB{
		Id:         1,
		CourierId:  1,
		OrderId:    "1",
		AssignedAt: time.Time{},
		Deadline:   time.Time{},
	}

	dRepo.EXPECT().Delete(gomock.Any(), orderId).Return(int64(1), nil)

	cRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, r *model.UpdateCourierRequest) error {
			require.NotNil(t, r.Id)
			require.NotNil(t, r.Status)
			require.Equal(t, int64(1), *r.Id)
			require.Equal(t, model.CourierStatusAvailable, *r.Status)
			return nil
		})

	expected := &model.UnassignCourier{
		CourierId: dModel.CourierId,
		OrderId:   orderId,
		Status:    model.Unassigned,
	}

	got, err := service.UnassignCourier(context.Background(), orderId)

	require.Equal(t, expected, got)
	require.NoError(t, err)

}

func TestUnassign_NotAssignedCourier(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"

	dRepo.EXPECT().Delete(gomock.Any(), orderId).Return(int64(0), delivery_repository.ErrNotFound)

	got, err := service.UnassignCourier(context.Background(), orderId)

	require.Nil(t, got)
	require.ErrorIs(t, ErrNotAssignedCourier, err)
}

func TestUnassign_DBError_GetByOrderId(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	dbErr := errors.New("db error")

	orderId := "1"

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	dRepo.EXPECT().Delete(gomock.Any(), orderId).Return(int64(0), dbErr)

	got, err := service.UnassignCourier(context.Background(), orderId)

	require.Nil(t, got)
	require.ErrorIs(t, dbErr, err)

}

func TestUnassign_DBError_Delete(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"
	dbErr := errors.New("db error")

	dRepo.EXPECT().Delete(gomock.Any(), orderId).Return(int64(0), dbErr)

	got, err := service.UnassignCourier(context.Background(), orderId)

	require.Nil(t, got)
	require.ErrorIs(t, dbErr, err)

}

func TestUnassign_DBError_Update(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"
	dbErr := errors.New("db error")

	dModel := &model.DeliveryDB{
		Id:         1,
		CourierId:  1,
		OrderId:    "1",
		AssignedAt: time.Time{},
		Deadline:   time.Time{},
	}

	dRepo.EXPECT().Delete(gomock.Any(), orderId).Return(int64(1), nil)

	var status model.CourierStatus

	status = model.CourierStatusAvailable

	req := &model.UpdateCourierRequest{
		Id:     &dModel.Id,
		Status: &status,
	}

	cRepo.EXPECT().Update(gomock.Any(), req).Return(dbErr)

	got, err := service.UnassignCourier(context.Background(), orderId)

	require.Nil(t, got)
	require.ErrorIs(t, dbErr, err)

}

func TestUnassign_DBError_Begin(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)
	orderID := "1"
	dbErr := errors.New("tx begin error")

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		Return(dbErr)

	got, err := service.UnassignCourier(context.Background(), orderID)

	require.Nil(t, got)
	require.Equal(t, dbErr, err)

}

func TestCompleteCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"
	delivery := &model.DeliveryDB{
		Id:        1,
		OrderId:   orderId,
		CourierId: 1,
	}

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderId).
		Return(delivery, nil)

	cRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, r *model.UpdateCourierRequest) error {
			require.NotNil(t, r.Id)
			require.NotNil(t, r.Status)
			require.Equal(t, int64(1), *r.Id)
			require.Equal(t, model.CourierStatusAvailable, *r.Status)
			return nil
		})

	err := service.CompleteCourier(context.Background(), orderId)
	require.NoError(t, err)
}

func TestCompleteCourier_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderId).
		Return(nil, delivery_repository.ErrNotFound)

	err := service.CompleteCourier(context.Background(), orderId)
	require.ErrorIs(t, err, ErrNotFoundOrder)
}

func TestCompleteCourier_DBError_GetByOrderId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"
	dbErr := errors.New("db error")

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderId).
		Return(nil, dbErr)

	err := service.CompleteCourier(context.Background(), orderId)
	require.ErrorIs(t, err, dbErr)
}

func TestCompleteCourier_DBError_UpdateCourier(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		DoAndReturn(func(parent context.Context, withTx bool, fn func(ctx context.Context) error) error {
			return fn(parent)
		})

	orderId := "1"
	delivery := &model.DeliveryDB{
		Id:        1,
		OrderId:   orderId,
		CourierId: 1,
	}

	dRepo.EXPECT().
		GetByOrderId(gomock.Any(), orderId).
		Return(delivery, nil)

	dbErr := errors.New("db error")

	cRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(dbErr)

	err := service.CompleteCourier(context.Background(), orderId)
	require.ErrorIs(t, err, dbErr)
}

func TestCompleteCourier_DBError_Begin(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cRepo := mocks.NewMockCourierRepository(ctrl)
	tx := mocks.NewMockTransactionManager(ctrl)
	dRepo := mocks.NewMockDeliveryRepository(ctrl)
	transportFactory := mocks.NewMockTransportFactory(ctrl)

	service := NewAssignService(tx, dRepo, cRepo, transportFactory)

	orderId := "1"
	dbErr := errors.New("tx begin error")

	tx.EXPECT().
		Begin(gomock.Any(), true, gomock.Any()).
		Return(dbErr)

	err := service.CompleteCourier(context.Background(), orderId)
	require.ErrorIs(t, err, dbErr)
}
