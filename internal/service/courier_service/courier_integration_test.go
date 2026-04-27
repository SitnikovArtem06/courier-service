//go:build integration

package courier_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/repository/courier_repository"
	"course-go-avito-SitnikovArtem06/internal/tx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"testing"
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

func truncateCouriers(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()
	_, err := pool.Exec(ctx, `TRUNCATE TABLE couriers RESTART IDENTITY CASCADE;`)
	require.NoError(t, err)
}

func newTestCourierService(t *testing.T) *CourierService {
	pool := newTestPool(t)
	truncateCouriers(t, pool)

	tm := tx.NewPgxTxManager(pool)
	repo := courier_repository.NewCourierRepo(tm)

	svc := NewCourierService(repo)

	return svc
}

func TestCreateCourier_Success_Integration(t *testing.T) {
	svc := newTestCourierService(t)
	ctx := context.Background()

	req := &model.CreateCourierRequest{
		Name:   "Courier",
		Phone:  "+70000000001",
		Status: model.CourierStatusAvailable,
	}

	c, err := svc.CreateCourier(ctx, req)
	require.NoError(t, err)

	require.NotZero(t, c.Id)
	require.Equal(t, "Courier", c.Name)
	require.Equal(t, "+70000000001", c.Phone)
	require.Equal(t, model.CourierStatusAvailable, c.Status)
	require.Equal(t, model.OnFoot, c.Transport)
}

func TestCreateCourier_DuplicatePhone_Integration(t *testing.T) {
	svc := newTestCourierService(t)
	ctx := context.Background()

	req := &model.CreateCourierRequest{
		Name:   "First",
		Phone:  "+70000000002",
		Status: model.CourierStatusAvailable,
	}

	_, err := svc.CreateCourier(ctx, req)
	require.NoError(t, err)

	req2 := &model.CreateCourierRequest{
		Name:   "Second",
		Phone:  "+70000000002",
		Status: model.CourierStatusAvailable,
	}

	_, err = svc.CreateCourier(ctx, req2)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDuplicatePhone)
}

func TestGetCourierById_Success_Integration(t *testing.T) {
	svc := newTestCourierService(t)
	ctx := context.Background()

	created, err := svc.CreateCourier(ctx, &model.CreateCourierRequest{
		Name:      "A",
		Phone:     "+70000000003",
		Status:    model.CourierStatusAvailable,
		Transport: model.OnFoot,
	})
	require.NoError(t, err)

	got, err := svc.GetCourierById(ctx, created.Id)
	require.NoError(t, err)

	require.Equal(t, created.Id, got.Id)
	require.Equal(t, created.Name, got.Name)
	require.Equal(t, created.Phone, got.Phone)
	require.Equal(t, created.Status, got.Status)
	require.Equal(t, created.Transport, got.Transport)
}

func TestGetCourierById_NotFound_Integration(t *testing.T) {
	svc := newTestCourierService(t)
	ctx := context.Background()

	_, err := svc.GetCourierById(ctx, 999999)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGetAllCouriers_Success_Integration(t *testing.T) {
	svc := newTestCourierService(t)
	ctx := context.Background()

	_, err := svc.CreateCourier(ctx, &model.CreateCourierRequest{
		Name:   "All 1",
		Phone:  "+70000000010",
		Status: model.CourierStatusAvailable,
	})
	require.NoError(t, err)

	_, err = svc.CreateCourier(ctx, &model.CreateCourierRequest{
		Name:   "All 2",
		Phone:  "+70000000011",
		Status: model.CourierStatusBusy,
	})
	require.NoError(t, err)

	list, err := svc.GetAllCouriers(ctx)
	require.NoError(t, err)

	require.Len(t, list, 2)
	require.Equal(t, "All 1", list[0].Name)
	require.Equal(t, "All 2", list[1].Name)
}

func TestUpdateCourier_Success(t *testing.T) {
	svc := newTestCourierService(t)
	ctx := context.Background()

	created, err := svc.CreateCourier(ctx, &model.CreateCourierRequest{
		Name:   "Old Name",
		Phone:  "+70000000020",
		Status: model.CourierStatusAvailable,
	})
	require.NoError(t, err)

	id := created.Id
	newName := "New Name"
	newPhone := "+70000000021"
	newStatus := model.CourierStatusAvailable
	newTransport := model.Car

	req := &model.UpdateCourierRequest{
		Id:        &id,
		Name:      &newName,
		Phone:     &newPhone,
		Status:    &newStatus,
		Transport: &newTransport,
	}

	err = svc.UpdateCourier(ctx, req)
	require.NoError(t, err)

	got, err := svc.GetCourierById(ctx, id)
	require.NoError(t, err)

	require.Equal(t, newName, got.Name)
	require.Equal(t, newPhone, got.Phone)
	require.Equal(t, newStatus, got.Status)
	require.Equal(t, newTransport, got.Transport)
}

func TestUpdateCourier_NotFound_Integration(t *testing.T) {
	svc := newTestCourierService(t)
	ctx := context.Background()

	id := int64(9)
	name := "A"

	req := &model.UpdateCourierRequest{
		Id:   &id,
		Name: &name,
	}

	err := svc.UpdateCourier(ctx, req)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestUpdateCourier_DuplicatePhone_Integration(t *testing.T) {
	svc := newTestCourierService(t)
	ctx := context.Background()

	first, err := svc.CreateCourier(ctx, &model.CreateCourierRequest{
		Name:   "First",
		Phone:  "+70000000030",
		Status: model.CourierStatusAvailable,
	})
	require.NoError(t, err)

	second, err := svc.CreateCourier(ctx, &model.CreateCourierRequest{
		Name:   "Second",
		Phone:  "+70000000031",
		Status: model.CourierStatusAvailable,
	})
	require.NoError(t, err)

	id := second.Id
	newPhone := first.Phone

	req := &model.UpdateCourierRequest{
		Id:    &id,
		Phone: &newPhone,
	}

	err = svc.UpdateCourier(ctx, req)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDuplicatePhone)
}
