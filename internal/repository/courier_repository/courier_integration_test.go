//go:build integration

package courier_repository

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
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

func truncateTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		TRUNCATE TABLE couriers, delivery RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}

func newTestCourierRepo(t *testing.T) *CourierRepo {
	pool := newTestPool(t)
	truncateTables(t, pool)

	tm := tx.NewPgxTxManager(pool)
	return NewCourierRepo(tm)
}

func TestCreateAndGet_Success_Integration(t *testing.T) {
	repo := newTestCourierRepo(t)
	ctx := context.Background()

	c := &model.CourierDB{
		Name:      "Courier",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		Transport: model.Car,
	}

	created, err := repo.Create(ctx, c)
	require.NoError(t, err)
	require.NotZero(t, created.Id)

	got, err := repo.Get(ctx, c.Id)
	require.NoError(t, err)

	require.Equal(t, c.Id, got.Id)
	require.Equal(t, c.Name, got.Name)
	require.Equal(t, c.Status, got.Status)
	require.Equal(t, c.Transport, got.Transport)
}

func TestCreate_DuplicatePhone_Integration(t *testing.T) {

	repo := newTestCourierRepo(t)
	ctx := context.Background()

	base := &model.CourierDB{
		Name:      "First",
		Phone:     "+70000000002",
		Status:    model.CourierStatusAvailable,
		Transport: model.OnFoot,
	}

	_, err := repo.Create(ctx, base)
	require.NoError(t, err)

	dup := &model.CourierDB{
		Name:      "Second",
		Phone:     base.Phone,
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	}

	_, err = repo.Create(ctx, dup)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDuplicatePhoneRepo)
}

func TestGet_NotFound_Integration(t *testing.T) {
	repo := newTestCourierRepo(t)
	ctx := context.Background()

	_, err := repo.Get(ctx, 999)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFoundRepo)
}

func TestGetAll_Success_Integration(t *testing.T) {
	repo := newTestCourierRepo(t)
	ctx := context.Background()

	c1, err := repo.Create(ctx, &model.CourierDB{
		Name:      "Courier 1",
		Phone:     "+70000000010",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	c2, err := repo.Create(ctx, &model.CourierDB{
		Name:      "Courier 2",
		Phone:     "+70000000011",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	couriers, err := repo.GetAll(ctx)
	require.NoError(t, err)

	require.Len(t, couriers, 2)
	require.Equal(t, c1.Id, couriers[0].Id)
	require.Equal(t, c2.Id, couriers[1].Id)
}

func TestUpdate_Success_Integration(t *testing.T) {
	repo := newTestCourierRepo(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, &model.CourierDB{
		Name:      "Old Name",
		Phone:     "+70000000020",
		Status:    model.CourierStatusAvailable,
		Transport: model.OnFoot,
	})
	require.NoError(t, err)

	id := created.Id
	newName := "New Name"
	newPhone := "+70000000021"
	newStatus := model.CourierStatusAvailable
	newTransport := model.OnFoot

	req := &model.UpdateCourierRequest{
		Id:        &id,
		Name:      &newName,
		Phone:     &newPhone,
		Status:    &newStatus,
		Transport: &newTransport,
	}

	err = repo.Update(ctx, req)
	require.NoError(t, err)

	got, err := repo.Get(ctx, id)
	require.NoError(t, err)

	require.Equal(t, newName, got.Name)
	require.Equal(t, newPhone, got.Phone)
	require.Equal(t, newStatus, got.Status)
	require.Equal(t, newTransport, got.Transport)
}

func TestUpdate_NotFound_Integration(t *testing.T) {
	repo := newTestCourierRepo(t)
	ctx := context.Background()

	id := int64(9)
	name := "A"

	req := &model.UpdateCourierRequest{
		Id:   &id,
		Name: &name,
	}

	err := repo.Update(ctx, req)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFoundRepo)
}

func TestUpdate_DuplicatePhone_Integration(t *testing.T) {
	repo := newTestCourierRepo(t)
	ctx := context.Background()

	c1, err := repo.Create(ctx, &model.CourierDB{
		Name:      "Courier 1",
		Phone:     "+70000000030",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	c2, err := repo.Create(ctx, &model.CourierDB{
		Name:      "Courier 2",
		Phone:     "+70000000031",
		Status:    model.CourierStatusAvailable,
		Transport: model.OnFoot,
	})
	require.NoError(t, err)

	id := c2.Id
	newPhone := c1.Phone

	req := &model.UpdateCourierRequest{
		Id:    &id,
		Phone: &newPhone,
	}

	err = repo.Update(ctx, req)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrDuplicatePhoneRepo)
}

func TestGetAvailableCouriers_Success_Integration(t *testing.T) {
	repo := newTestCourierRepo(t)
	ctx := context.Background()

	c1, err := repo.Create(ctx, &model.CourierDB{
		Name:      "Avail 1",
		Phone:     "+70000000040",
		Status:    model.CourierStatusAvailable,
		Transport: model.Car,
	})
	require.NoError(t, err)

	_, err = repo.Create(ctx, &model.CourierDB{
		Name:      "Busy 1",
		Phone:     "+70000000041",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	c3, err := repo.Create(ctx, &model.CourierDB{
		Name:      "Avail 2",
		Phone:     "+70000000042",
		Status:    model.CourierStatusAvailable,
		Transport: model.Car,
	})
	require.NoError(t, err)

	got, err := repo.GetAvailableCouriers(ctx)
	require.NoError(t, err)

	require.Len(t, got, 2)
	require.Equal(t, c1.Id, got[0].Id)
	require.Equal(t, c3.Id, got[1].Id)
	for _, c := range got {
		require.Equal(t, model.CourierStatusAvailable, c.Status)
	}
}

func TestUpdateAllExpiredCourier_Success_Integration(t *testing.T) {
	repo := newTestCourierRepo(t)
	ctx := context.Background()

	busy1, err := repo.Create(ctx, &model.CourierDB{
		Name:      "Busy 1",
		Phone:     "+70000000050",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	busy2, err := repo.Create(ctx, &model.CourierDB{
		Name:      "Busy 2",
		Phone:     "+70000000051",
		Status:    model.CourierStatusBusy,
		Transport: model.Car,
	})
	require.NoError(t, err)

	avail, err := repo.Create(ctx, &model.CourierDB{
		Name:      "Avail",
		Phone:     "+70000000052",
		Status:    model.CourierStatusAvailable,
		Transport: model.OnFoot,
	})
	require.NoError(t, err)

	ids := []int64{busy1.Id, busy2.Id, avail.Id}

	err = repo.UpdateAllExpiredCourier(ctx, ids)
	require.NoError(t, err)

	b1, err := repo.Get(ctx, busy1.Id)
	require.NoError(t, err)
	b2, err := repo.Get(ctx, busy2.Id)
	require.NoError(t, err)
	a, err := repo.Get(ctx, avail.Id)
	require.NoError(t, err)

	require.Equal(t, model.CourierStatusAvailable, b1.Status)
	require.Equal(t, model.CourierStatusAvailable, b2.Status)
	require.Equal(t, model.CourierStatusAvailable, a.Status)
}
