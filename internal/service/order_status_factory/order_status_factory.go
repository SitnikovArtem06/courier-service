package order_status_factory

import "context"

type OrderStatus interface {
	Do(ctx context.Context, orderId string) error
}

type OrderStatusImpl struct {
	created   OrderStatus
	cancelled OrderStatus
	completed OrderStatus
}

func NewOrderStatusFactory(a assign) *OrderStatusImpl {
	return &OrderStatusImpl{
		created:   &Created{s: a},
		cancelled: &Cancelled{s: a},
		completed: &Completed{s: a},
	}
}

type OrderStatusFactory interface {
	Get(status string) OrderStatus
}

func (f *OrderStatusImpl) Get(status string) OrderStatus {
	switch status {
	case "created":
		return f.created
	case "cancelled":
		return f.cancelled
	case "completed":
		return f.completed
	default:
		return nil
	}
}

type Created struct {
	s assign
}

func (c Created) Do(ctx context.Context, orderId string) error {

	_, err := c.s.AssignCourier(ctx, orderId)
	return err
}

type Cancelled struct {
	s assign
}

func (c Cancelled) Do(ctx context.Context, orderId string) error {

	_, err := c.s.UnassignCourier(ctx, orderId)
	return err
}

type Completed struct {
	s assign
}

func (c Completed) Do(ctx context.Context, orderId string) error {
	return c.s.CompleteCourier(ctx, orderId)
}
