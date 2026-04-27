package assign_handler

import (
	"bytes"
	assign_handler "course-go-avito-SitnikovArtem06/internal/handlers/assign_handler/mocks"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/service/assign_service"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAssignCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := assign_handler.NewMockassignService(ctrl)

	h := NewAssignHandler(svc)

	orderID := "123"
	deadline := time.Now().UTC()

	assignModel := &model.AssignCourier{
		CourierId: 1,
		OrderId:   orderID,
		Transport: model.Car,
		Deadline:  deadline,
	}

	svc.EXPECT().
		AssignCourier(gomock.Any(), orderID).
		Return(assignModel, nil)

	body, _ := json.Marshal(map[string]string{
		"order_id": orderID,
	})

	req := httptest.NewRequest(http.MethodPost, "/delivery/assign", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.AssignCourier(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp assignCourierResp
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, int64(1), resp.CourierId)
	require.Equal(t, orderID, resp.OrderId)
	require.Equal(t, model.Car.String(), resp.Transport)
	require.WithinDuration(t, deadline, resp.Deadline, time.Second)
}

func TestAssignCourier_InvalidOrderID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := assign_handler.NewMockassignService(ctrl)
	h := NewAssignHandler(svc)

	body, _ := json.Marshal(map[string]string{
		"order_id": "",
	})

	req := httptest.NewRequest(http.MethodPost, "/delivery/assign", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.AssignCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "invalid order_id", resp["error"])
}

func TestAssignCourier_NotAvailable(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := assign_handler.NewMockassignService(ctrl)
	h := NewAssignHandler(svc)

	orderID := "123"

	svc.EXPECT().
		AssignCourier(gomock.Any(), orderID).
		Return(nil, assign_service.ErrNotAvailableCourier)

	body, _ := json.Marshal(map[string]string{
		"order_id": orderID,
	})

	req := httptest.NewRequest(http.MethodPost, "/delivery/assign", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.AssignCourier(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, assign_service.ErrNotAvailableCourier.Error(), resp["error"])
}

func TestAssignCourier_AlreadyAssigned(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := assign_handler.NewMockassignService(ctrl)
	h := NewAssignHandler(svc)

	orderID := "123"

	svc.EXPECT().
		AssignCourier(gomock.Any(), orderID).
		Return(nil, assign_service.ErrOrderAlreadyAssign)

	body, _ := json.Marshal(map[string]string{
		"order_id": orderID,
	})

	req := httptest.NewRequest(http.MethodPost, "/delivery/assign", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.AssignCourier(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, assign_service.ErrOrderAlreadyAssign.Error(), resp["error"])
}

func TestAssignCourier_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := assign_handler.NewMockassignService(ctrl)
	h := NewAssignHandler(svc)

	orderID := "123"
	internalErr := errors.New("db error")

	svc.EXPECT().
		AssignCourier(gomock.Any(), orderID).
		Return(nil, internalErr)

	body, _ := json.Marshal(map[string]string{
		"order_id": orderID,
	})

	req := httptest.NewRequest(http.MethodPost, "/delivery/assign", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.AssignCourier(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUnassignCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := assign_handler.NewMockassignService(ctrl)
	h := NewAssignHandler(svc)

	orderID := "123"

	svc.EXPECT().
		UnassignCourier(gomock.Any(), orderID).
		Return(&model.UnassignCourier{
			CourierId: 1,
			OrderId:   orderID,
			Status:    model.Unassigned,
		}, nil)

	body, _ := json.Marshal(map[string]string{
		"order_id": orderID,
	})

	req := httptest.NewRequest(http.MethodPost, "/delivery/unassign", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UnassignCourier(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp unassignCourierResp
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	require.Equal(t, int64(1), resp.CourierId)
	require.Equal(t, orderID, resp.OrderId)
	require.Equal(t, model.Unassigned.String(), resp.Status)
}

func TestUnassignCourier_InvalidOrderID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := assign_handler.NewMockassignService(ctrl)
	h := NewAssignHandler(svc)

	body, _ := json.Marshal(map[string]string{
		"order_id": "",
	})

	req := httptest.NewRequest(http.MethodPost, "/delivery/unassign", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UnassignCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "invalid order_id", resp["error"])
}

func TestUnassignCourier_NotAssigned(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := assign_handler.NewMockassignService(ctrl)
	h := NewAssignHandler(svc)

	orderID := "123"

	svc.EXPECT().
		UnassignCourier(gomock.Any(), orderID).
		Return(nil, assign_service.ErrNotAssignedCourier)

	body, _ := json.Marshal(map[string]string{
		"order_id": orderID,
	})

	req := httptest.NewRequest(http.MethodPost, "/delivery/unassign", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UnassignCourier(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, assign_service.ErrNotAssignedCourier.Error(), resp["error"])
}

func TestUnassignCourier_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := assign_handler.NewMockassignService(ctrl)
	h := NewAssignHandler(svc)

	orderID := "123"
	internalErr := errors.New("db error")

	svc.EXPECT().
		UnassignCourier(gomock.Any(), orderID).
		Return(nil, internalErr)

	body, _ := json.Marshal(map[string]string{
		"order_id": orderID,
	})

	req := httptest.NewRequest(http.MethodPost, "/delivery/unassign", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UnassignCourier(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
