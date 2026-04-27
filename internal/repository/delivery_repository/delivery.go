package delivery_repository

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/tx"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type DeliveryRepo struct {
	tm tx.TransactionManager
}

func NewDeliveryRepository(tm tx.TransactionManager) *DeliveryRepo {
	return &DeliveryRepo{tm: tm}
}

func (r *DeliveryRepo) Create(ctx context.Context, orderId string, courierId int64, deadline time.Time) error {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return err
	}

	sqlInsert := `INSERT INTO delivery (courier_id,order_id, deadline) VALUES($1,$2,$3)`

	if _, err := conn.Exec(ctx, sqlInsert, courierId, orderId, deadline); err != nil {
		return err
	}

	return nil

}

func (r *DeliveryRepo) GetByOrderId(ctx context.Context, orderID string) (*model.DeliveryDB, error) {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return nil, err
	}

	sqlSelect := `SELECT id, courier_id, order_id, assigned_at, deadline FROM delivery WHERE order_id = $1 FOR UPDATE;`

	var delivery model.DeliveryDB

	if err := conn.QueryRow(ctx, sqlSelect, orderID).Scan(&delivery.Id, &delivery.CourierId, &delivery.OrderId, &delivery.AssignedAt, &delivery.Deadline); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &delivery, nil
}

func (r *DeliveryRepo) Delete(ctx context.Context, orderId string) (int64, error) {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return 0, err
	}

	sqlDelete := `DELETE FROM delivery WHERE order_id = $1 RETURNING courier_id;`

	var courierId int64
	err = conn.QueryRow(ctx, sqlDelete, orderId).Scan(&courierId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	return courierId, nil

}

func (r *DeliveryRepo) GetExpiredOrders(ctx context.Context) ([]int64, error) {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return nil, err
	}

	sqlSelect := `SELECT DISTINCT d.courier_id FROM delivery d WHERE d.deadline < now() AND NOT EXISTS (SELECT 1
      				FROM delivery d2
      				WHERE d2.courier_id = d.courier_id
        			AND d2.deadline >= now()
 		 );`

	rows, err := conn.Query(ctx, sqlSelect)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var id int64

	courier_ids := make([]int64, 0)

	for rows.Next() {

		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		courier_ids = append(courier_ids, id)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courier_ids, nil

}

func (r *DeliveryRepo) GetCourierWithMinimumOrder(ctx context.Context) (id int64, err error) {

	conn, err := r.tm.GetConnection(ctx)
	if err != nil {
		return 0, err
	}

	sqlSelect := `SELECT c.id FROM couriers c LEFT JOIN delivery d ON c.id = d.courier_id AND d.deadline < now() WHERE c.status = 'available' GROUP BY c.id ORDER BY COUNT(d.order_id) ASC, c.id ASC LIMIT 1;`

	var courierId int64

	if err = conn.QueryRow(ctx, sqlSelect).Scan(&courierId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoneCourier
		}
		return 0, err
	}

	return courierId, nil

}
