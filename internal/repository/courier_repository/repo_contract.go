package courier_repository

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
)

type CourierRepository interface {
	Create(ctx context.Context, courier *model.CourierDB) (*model.CourierDB, error)

	Get(ctx context.Context, id int64) (*model.CourierDB, error)

	GetAll(ctx context.Context) ([]model.CourierDB, error)

	Update(ctx context.Context, req *model.UpdateCourierRequest) error

	GetAvailableCouriers(ctx context.Context) ([]model.CourierDB, error)
	UpdateAllExpiredCourier(ctx context.Context, ids []int64) error
}
