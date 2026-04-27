package assign_handler

import (
	"course-go-avito-SitnikovArtem06/internal/service/assign_service"
	"encoding/json"
	"errors"
	"net/http"
)

type AssignHandler struct {
	as assignService
}

func NewAssignHandler(service assignService) *AssignHandler {
	return &AssignHandler{as: service}
}

func (h *AssignHandler) AssignCourier(w http.ResponseWriter, r *http.Request) {

	var orderReq order

	json.NewDecoder(r.Body).Decode(&orderReq)

	if orderReq.OrderId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid order_id",
		})
		return
	}

	assign, err := h.as.AssignCourier(r.Context(), orderReq.OrderId)
	if err != nil {
		switch {
		case errors.Is(err, assign_service.ErrNotAvailableCourier):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		case errors.Is(err, assign_service.ErrOrderAlreadyAssign):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return

	}

	resp := assignCourierResp{
		CourierId: assign.CourierId,
		OrderId:   assign.OrderId,
		Transport: assign.Transport.String(),
		Deadline:  assign.Deadline,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func (h *AssignHandler) UnassignCourier(w http.ResponseWriter, r *http.Request) {

	var orderReq order

	json.NewDecoder(r.Body).Decode(&orderReq)

	if orderReq.OrderId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid order_id",
		})
		return
	}

	unassign, err := h.as.UnassignCourier(r.Context(), orderReq.OrderId)
	if err != nil {
		switch {
		case errors.Is(err, assign_service.ErrNotAssignedCourier):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	resp := unassignCourierResp{
		OrderId:   unassign.OrderId,
		Status:    unassign.Status.String(),
		CourierId: unassign.CourierId,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
