//go:build integration

package assign_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/repository/courier_repository"
	"course-go-avito-SitnikovArtem06/internal/repository/delivery_repository"
	"course-go-avito-SitnikovArtem06/internal/service/transport_factory"
	"course-go-avito-SitnikovArtem06/internal/tx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := "postgres://sito:sito@localhost:5432/test_db"

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	t.Cleanup(func() { pool.Close() })
	return pool
}

func truncateTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE delivery, couriers RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}

func newTestAssignService(t *testing.T) (*AssignService, *courier_repository.CourierRepo, *delivery_repository.DeliveryRepo) {
	pool := newTestPool(t)
	truncateTables(t, pool)

	tm := tx.NewPgxTxManager(pool)
	cRepo := courier_repository.NewCourierRepo(tm)
	dRepo := delivery_repository.NewDeliveryRepository(tm)
	tf := transport_factory.NewTransportFactory()

	svc := NewAssignService(tm, dRepo, cRepo, tf)

	return svc, cRepo, dRepo
}

func TestAssignCourier_Success_Integration(t *testing.T) {
	svc, cRepo, dRepo := newTestAssignService(t)
	ctx := context.Background()

	_, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "Courier 1",
		Phone:     "+79990000001",
		Status:    model.CourierStatusAvailable,
		Transport: model.Car,
	})
	require.NoError(t, err)

	_, err = cRepo.Create(ctx, &model.CourierDB{
		Name:      "Courier 2",
		Phone:     "+79990000002",
		Status:    model.CourierStatusAvailable,
		Transport: model.OnFoot,
	})
	require.NoError(t, err)

	orderID := "order-1"

	res, err := svc.AssignCourier(ctx, orderID)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, orderID, res.OrderId)
	require.NotZero(t, res.CourierId)

	require.True(t, res.Deadline.After(time.Now().Add(-1*time.Second)))

	d, err := dRepo.GetByOrderId(ctx, orderID)
	require.NoError(t, err)
	require.Equal(t, res.CourierId, d.CourierId)
	require.Equal(t, orderID, d.OrderId)

	courier, err := cRepo.Get(ctx, res.CourierId)
	require.NoError(t, err)
	require.Equal(t, model.CourierStatusBusy, courier.Status)
}

func TestAssignCourier_NoAvailableCouriers_Integration(t *testing.T) {
	svc, cRepo, _ := newTestAssignService(t)
	ctx := context.Background()

	_, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "Busy 1",
		Phone:     "+79990000003",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	orderID := "order-no-available"

	res, err := svc.AssignCourier(ctx, orderID)
	require.Nil(t, res)
	require.ErrorIs(t, err, ErrNotAvailableCourier)
}

func TestAssignCourier_AlreadyAssigned_Integration(t *testing.T) {
	svc, cRepo, dRepo := newTestAssignService(t)
	ctx := context.Background()

	courier, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "Already",
		Phone:     "+79990000004",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	orderID := "order-already"

	deadline := time.Now().Add(30 * time.Minute)
	err = dRepo.Create(ctx, orderID, courier.Id, deadline)
	require.NoError(t, err)

	res, err := svc.AssignCourier(ctx, orderID)
	require.Nil(t, res)
	require.ErrorIs(t, err, ErrOrderAlreadyAssign)
}

func TestUnassignCourier_Success_Integration(t *testing.T) {
	svc, cRepo, dRepo := newTestAssignService(t)
	ctx := context.Background()

	courier, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "Busy",
		Phone:     "+79990000005",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	orderID := "order-unassign-1"
	deadline := time.Now().Add(30 * time.Minute)

	err = dRepo.Create(ctx, orderID, courier.Id, deadline)
	require.NoError(t, err)

	res, err := svc.UnassignCourier(ctx, orderID)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, orderID, res.OrderId)
	require.Equal(t, courier.Id, res.CourierId)
	require.Equal(t, model.Unassigned, res.Status)

	_, err = dRepo.GetByOrderId(ctx, orderID)
	require.Error(t, err)

	updated, err := cRepo.Get(ctx, courier.Id)
	require.NoError(t, err)
	require.Equal(t, model.CourierStatusAvailable, updated.Status)
}

func TestUnassignCourier_NotAssigned_Integration(t *testing.T) {
	svc, _, _ := newTestAssignService(t)
	ctx := context.Background()

	orderID := "order-not-assigned"

	res, err := svc.UnassignCourier(ctx, orderID)
	require.Nil(t, res)
	require.ErrorIs(t, err, ErrNotAssignedCourier)
}
