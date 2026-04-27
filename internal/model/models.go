package model

import (
	"time"
)

type Courier struct {
	Id        int64
	Name      string
	Phone     string
	Status    CourierStatus
	Transport TransportType
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateCourierRequest struct {
	Name      string
	Phone     string
	Status    CourierStatus
	Transport TransportType
}

type UpdateCourierRequest struct {
	Id        *int64
	Name      *string
	Phone     *string
	Status    *CourierStatus
	Transport *TransportType
}

type AssignCourier struct {
	CourierId int64
	OrderId   string
	Transport TransportType
	Deadline  time.Time
}

type UnassignCourier struct {
	CourierId int64
	OrderId   string
	Status    AssignStatus
}

type OrdersResponse struct {
	OrdersId  []string
	CreatedAt []time.Time
}

type ChangedStatus struct {
	OrderID string
	Status  string
}
type CourierStatus string

const (
	CourierStatusAvailable CourierStatus = "available"
	CourierStatusBusy      CourierStatus = "busy"
	CourierStatusPaused    CourierStatus = "paused"
)

func (s CourierStatus) IsValid() bool {
	return s == CourierStatusAvailable || s == CourierStatusBusy || s == CourierStatusPaused
}

func (s CourierStatus) String() string {
	return string(s)
}

type TransportType string

const (
	OnFoot  TransportType = "on_foot"
	Scooter TransportType = "scooter"

	Car TransportType = "car"
)

func (t TransportType) IsValid() bool {
	return t == OnFoot || t == Scooter || t == Car
}

func (t TransportType) String() string {
	return string(t)
}

type AssignStatus string

const (
	Assigned AssignStatus = "assigned"

	Unassigned AssignStatus = "unassigned"
)

func (a AssignStatus) String() string {
	return string(a)
}
