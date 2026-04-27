package model

import (
	"time"
)

type CourierDB struct {
	Id        int64         `db:"id"`
	Name      string        `db:"name"`
	Phone     string        `db:"phone"`
	Status    CourierStatus `db:"status"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	Transport TransportType `db:"transport_type"`
}

type DeliveryDB struct {
	Id         int64     `db:"id"`
	CourierId  int64     `db:"courier_id"`
	OrderId    string    `db:"order_id"`
	AssignedAt time.Time `db:"assigned_at"`
	Deadline   time.Time `db:"deadline"`
}
