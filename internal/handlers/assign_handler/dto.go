package assign_handler

import "time"

type order struct {
	OrderId string `json:"order_id"`
}

type assignCourierResp struct {
	CourierId int64 `json:"courier_id"`

	OrderId string `json:"order_id"`

	Transport string `json:"transport_type"`

	Deadline time.Time `json:"delivery_deadline"`
}

type unassignCourierResp struct {
	OrderId string `json:"order_id"`

	Status string `json:"status"`

	CourierId int64 `json:"courier_id"`
}
