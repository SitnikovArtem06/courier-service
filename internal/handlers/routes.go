package handlers

import (
	"course-go-avito-SitnikovArtem06/internal/handlers/assign_handler"
	"course-go-avito-SitnikovArtem06/internal/handlers/courier_handler"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Routes(h *courier_handler.Handler, ha *assign_handler.AssignHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/courier/{id}", h.GetById)
	r.Get("/couriers", h.GetAll)
	r.Post("/courier", h.CreateCourier)
	r.Put("/courier", h.UpdateCourier)

	r.Post("/delivery/assign", ha.AssignCourier)
	r.Post("/delivery/unassign", ha.UnassignCourier)

	r.Get("/ping", Ping)
	r.Head("/healthcheck", HealthCheck)

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	return r
}
