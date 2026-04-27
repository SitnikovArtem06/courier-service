package courier_handler

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
)

type courierService interface {
	CreateCourier(ctx context.Context, c *model.CreateCourierRequest) (*model.Courier, error)
	GetCourierById(ctx context.Context, id int64) (*model.Courier, error)

	GetAllCouriers(ctx context.Context) ([]model.Courier, error)

	UpdateCourier(ctx context.Context, req *model.UpdateCourierRequest) error
}
