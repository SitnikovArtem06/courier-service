package changed

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/service/order_changed_service"
	"encoding/json"
	"errors"
	"fmt"
)

type ChangedHandler struct {
	changedS orderChanged
}

func NewChangedHandler(changedS orderChanged) *ChangedHandler {
	return &ChangedHandler{changedS: changedS}
}

func (h *ChangedHandler) HandleMessage(ctx context.Context, value []byte) error {
	var req OrderStatusChanged
	if err := json.Unmarshal(value, &req); err != nil {
		return fmt.Errorf("bad json: %v value=%s", err, string(value))
	}

	chg := model.ChangedStatus{
		OrderID: req.OrderID,
		Status:  req.Status,
	}

	if err := h.changedS.HandleStatusChanged(ctx, chg); err != nil {
		if errors.Is(err, order_changed_service.ErrMismatchStatus) {
			return nil
		}
		return err
	}
	return nil
}
