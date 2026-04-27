package courier_handler

import (
	"course-go-avito-SitnikovArtem06/internal/service/courier_service"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Handler struct {
	sc courierService
}

func NewHandler(sc courierService) *Handler {
	return &Handler{sc: sc}
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": ErrInvalidId.Error(),
		})
		return
	}

	if id <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": ErrInvalidId.Error(),
		})
		return
	}

	c, err := h.sc.GetCourierById(r.Context(), int64(id))

	if err != nil {

		switch {

		case errors.Is(err, courier_service.ErrNotFound):
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

	resp := toDTO(c)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(resp)

}

func (h *Handler) CreateCourier(w http.ResponseWriter, r *http.Request) {

	var reqDto createCourierDTO

	json.NewDecoder(r.Body).Decode(&reqDto)

	if err := reqDto.validateCreate(); err != nil {
		switch {

		case errors.Is(err, ErrEmptyName):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		}
		return
	}

	req := fromCreateDTO(reqDto)

	c, err := h.sc.CreateCourier(r.Context(), &req)

	if err != nil {
		switch {

		case errors.Is(err, courier_service.ErrInvalidStatus):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, courier_service.ErrInvalidPhoneNumber):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, courier_service.ErrDuplicatePhone):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, courier_service.ErrInvalidTransport):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{
		"id": c.Id,
	})

}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

	couriers, err := h.sc.GetAllCouriers(r.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)

	respCouriers := toDTOs(couriers)

	enc.SetIndent("", "  ")
	enc.Encode(respCouriers)

}

func (h *Handler) UpdateCourier(w http.ResponseWriter, r *http.Request) {

	var reqDto updateCourierDTO

	json.NewDecoder(r.Body).Decode(&reqDto)

	if err := reqDto.validateUpdate(); err != nil {
		switch {

		case errors.Is(err, ErrEmptyName):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, ErrInvalidId):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		}
		return
	}

	req := fromUpdateDTO(reqDto)

	err := h.sc.UpdateCourier(r.Context(), &req)

	if err != nil {

		switch {

		case errors.Is(err, courier_service.ErrInvalidStatus):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, courier_service.ErrInvalidPhoneNumber):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, courier_service.ErrNotFound):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, courier_service.ErrInvalidTransport):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, courier_service.ErrDuplicatePhone):
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

	w.WriteHeader(http.StatusOK)

}
