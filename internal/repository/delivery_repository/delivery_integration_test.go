//go:build integration

package delivery_repository

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/repository/courier_repository"
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

func newTestRepos(t *testing.T) (*DeliveryRepo, *courier_repository.CourierRepo) {
	pool := newTestPool(t)
	truncateTables(t, pool)

	tm := tx.NewPgxTxManager(pool)

	dRepo := NewDeliveryRepository(tm)
	cRepo := courier_repository.NewCourierRepo(tm)

	return dRepo, cRepo
}

func TestCreateAndGetByOrderId_Success_Integration(t *testing.T) {
	dRepo, cRepo := newTestRepos(t)
	ctx := context.Background()

	courier, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "Courier",
		Phone:     "+79990000001",
		Status:    model.CourierStatusAvailable,
		Transport: model.OnFoot,
	})
	require.NoError(t, err)

	orderID := "order-1"
	deadline := time.Now().Add(30 * time.Minute).UTC()

	err = dRepo.Create(ctx, orderID, courier.Id, deadline)
	require.NoError(t, err)

	got, err := dRepo.GetByOrderId(ctx, orderID)
	require.NoError(t, err)

	require.Equal(t, orderID, got.OrderId)
	require.Equal(t, courier.Id, got.CourierId)
	require.WithinDuration(t, deadline, got.Deadline, time.Second*2)
}

func TestGetByOrderId_NotFound_Integration(t *testing.T) {
	dRepo, _ := newTestRepos(t)
	ctx := context.Background()

	_, err := dRepo.GetByOrderId(ctx, "unknown-order")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestDelete_Success_Integration(t *testing.T) {
	dRepo, cRepo := newTestRepos(t)
	ctx := context.Background()

	courier, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "Courier",
		Phone:     "+79990000002",
		Status:    model.CourierStatusAvailable,
		Transport: model.OnFoot,
	})
	require.NoError(t, err)

	orderID := "order-delete"
	deadline := time.Now().Add(30 * time.Minute).UTC()

	err = dRepo.Create(ctx, orderID, courier.Id, deadline)
	require.NoError(t, err)

	deletedCourierId, err := dRepo.Delete(ctx, orderID)
	require.NoError(t, err)
	require.Equal(t, courier.Id, deletedCourierId)

	_, err = dRepo.GetByOrderId(ctx, orderID)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestDelete_NotFound_Integration(t *testing.T) {
	dRepo, _ := newTestRepos(t)
	ctx := context.Background()

	_, err := dRepo.Delete(ctx, "unknown-order")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGetExpiredOrders_Success_Integration(t *testing.T) {
	dRepo, cRepo := newTestRepos(t)
	ctx := context.Background()

	c1, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "ExpOnly",
		Phone:     "+79990000010",
		Status:    model.CourierStatusBusy,
		Transport: model.OnFoot,
	})
	require.NoError(t, err)

	c2, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "ExpAndFuture",
		Phone:     "+79990000011",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	c3, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "FutureOnly",
		Phone:     "+79990000012",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	past := time.Now().Add(-2 * time.Hour).UTC()
	future := time.Now().Add(2 * time.Hour).UTC()

	require.NoError(t, dRepo.Create(ctx, "order-c1-exp-1", c1.Id, past))

	require.NoError(t, dRepo.Create(ctx, "order-c2-exp-1", c2.Id, past))
	require.NoError(t, dRepo.Create(ctx, "order-c2-fut-1", c2.Id, future))

	require.NoError(t, dRepo.Create(ctx, "order-c3-fut-1", c3.Id, future))

	ids, err := dRepo.GetExpiredOrders(ctx)
	require.NoError(t, err)

	require.Len(t, ids, 1)
	require.Equal(t, c1.Id, ids[0])
}

func TestGetExpiredOrders_None_Integration(t *testing.T) {
	dRepo, cRepo := newTestRepos(t)
	ctx := context.Background()

	c, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "FutureOnly",
		Phone:     "+79990000013",
		Status:    model.CourierStatusBusy,
		Transport: model.OnFoot,
	})
	require.NoError(t, err)

	future := time.Now().Add(2 * time.Hour).UTC()

	require.NoError(t, dRepo.Create(ctx, "order-future-1", c.Id, future))

	ids, err := dRepo.GetExpiredOrders(ctx)
	require.NoError(t, err)
	require.Len(t, ids, 0)
}

func TestGetCourierWithMinimumOrder_Success_Integration(t *testing.T) {
	dRepo, cRepo := newTestRepos(t)
	ctx := context.Background()

	c1, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "ZeroExpired",
		Phone:     "+79990000020",
		Status:    model.CourierStatusAvailable,
		Transport: model.OnFoot,
	})
	require.NoError(t, err)

	c2, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "OneExpired",
		Phone:     "+79990000021",
		Status:    model.CourierStatusAvailable,
		Transport: model.Car,
	})
	require.NoError(t, err)

	c3, err := cRepo.Create(ctx, &model.CourierDB{
		Name:      "TwoExpired",
		Phone:     "+79990000022",
		Status:    model.CourierStatusAvailable,
		Transport: model.Car,
	})
	require.NoError(t, err)

	past := time.Now().Add(-2 * time.Hour).UTC()

	require.NoError(t, dRepo.Create(ctx, "order-c2-exp-1", c2.Id, past))

	require.NoError(t, dRepo.Create(ctx, "order-c3-exp-1", c3.Id, past))
	require.NoError(t, dRepo.Create(ctx, "order-c3-exp-2", c3.Id, past))

	gotID, err := dRepo.GetCourierWithMinimumOrder(ctx)
	require.NoError(t, err)
	require.Equal(t, c1.Id, gotID)
}

func TestGetCourierWithMinimumOrder_NoneCourier_Integration(t *testing.T) {
	dRepo, _ := newTestRepos(t)
	ctx := context.Background()

	_, err := dRepo.GetCourierWithMinimumOrder(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoneCourier)
}
