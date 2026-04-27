package changed

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
)

type orderChanged interface {
	HandleStatusChanged(ctx context.Context, req model.ChangedStatus) error
}
