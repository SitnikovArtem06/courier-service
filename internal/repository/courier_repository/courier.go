package courier_repository

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/tx"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"
)

type CourierRepo struct {
	tm tx.TransactionManager
}

func NewCourierRepo(tm tx.TransactionManager) *CourierRepo {
	return &CourierRepo{tm: tm}
}

func (r *CourierRepo) Create(ctx context.Context, courier *model.CourierDB) (*model.CourierDB, error) {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return nil, err
	}

	var id int64

	err = conn.QueryRow(ctx, `INSERT INTO couriers (name, phone,status,transport_type)
		VALUES ($1, $2, $3, $4) RETURNING id;`, courier.Name, courier.Phone, courier.Status, courier.Transport).Scan(&id)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, ErrDuplicatePhoneRepo
		}
		return nil, fmt.Errorf("database: %w", err)
	}

	courier.Id = id

	return courier, err

}

func (r *CourierRepo) Get(ctx context.Context, id int64) (*model.CourierDB, error) {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return nil, err
	}

	var courier model.CourierDB

	err = conn.QueryRow(ctx, `SELECT id, name, phone, status,created_at, updated_at, transport_type FROM couriers WHERE id=$1;`, id).Scan(&courier.Id, &courier.Name, &courier.Phone, &courier.Status, &courier.CreatedAt, &courier.UpdatedAt, &courier.Transport)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFoundRepo
		}
		return nil, fmt.Errorf("database: %w", err)
	}

	return &courier, nil

}

func (r *CourierRepo) GetAll(ctx context.Context) ([]model.CourierDB, error) {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, `SELECT id, name, phone, status, created_at, updated_at, transport_type FROM couriers ORDER BY id;`)

	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	defer rows.Close()

	var couriers []model.CourierDB

	for rows.Next() {
		var c model.CourierDB

		if err := rows.Scan(&c.Id, &c.Name, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt, &c.Transport); err != nil {
			return nil, err
		}

		couriers = append(couriers, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return couriers, nil

}

func (r *CourierRepo) Update(ctx context.Context, in *model.UpdateCourierRequest) error {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return err
	}

	tag, err := conn.Exec(ctx, `
        UPDATE couriers SET
            name       = COALESCE($2, name),
            phone      = COALESCE($3, phone),
            status     = COALESCE($4, status),
            transport_type = COALESCE($5, transport_type),
            updated_at = now()
        WHERE id = $1`,
		*in.Id, in.Name, in.Phone, in.Status, in.Transport)
	if err != nil {

		if strings.Contains(err.Error(), "duplicate key value") {
			return ErrDuplicatePhoneRepo
		}
		return fmt.Errorf("database: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFoundRepo
	}

	return nil
}

func (r *CourierRepo) GetAvailableCouriers(ctx context.Context) ([]model.CourierDB, error) {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return nil, err
	}

	sqlSelect := `SELECT id, name, phone, status, created_at, updated_at, transport_type FROM couriers WHERE status = 'available' ORDER BY id FOR UPDATE;`

	rows, err := conn.Query(ctx, sqlSelect)

	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	defer rows.Close()

	var couriers []model.CourierDB

	for rows.Next() {
		var c model.CourierDB

		if err := rows.Scan(&c.Id, &c.Name, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt, &c.Transport); err != nil {
			return nil, err
		}

		couriers = append(couriers, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return couriers, nil

}

func (r *CourierRepo) UpdateAllExpiredCourier(ctx context.Context, ids []int64) error {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return err
	}

	sqlUpdate := `UPDATE couriers SET status = 'available',
                    updated_at = now()
                    WHERE id = ANY($1) and status = 'busy'`

	if _, err = conn.Exec(ctx, sqlUpdate, ids); err != nil {
		return err
	}

	return nil

}
