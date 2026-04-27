package courier_handler

import (
	"bytes"
	"context"
	courier_handler "course-go-avito-SitnikovArtem06/internal/handlers/courier_handler/mocks"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/service/courier_service"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func withIDParam(req *http.Request, id string) *http.Request {
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", id)

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rc)
	return req.WithContext(ctx)
}

func TestGetById_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	courier := &model.Courier{
		Id:        1,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		Transport: model.OnFoot,
	}

	svc.EXPECT().
		GetCourierById(gomock.Any(), int64(1)).
		Return(courier, nil)

	req := httptest.NewRequest(http.MethodGet, "/courier/1", nil)
	req = withIDParam(req, "1")

	rec := httptest.NewRecorder()

	h.GetById(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp courierDTO
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	require.Equal(t, int64(1), resp.ID)
	require.Equal(t, "Artem", resp.Name)
	require.Equal(t, "+79119568101", resp.Phone)
	require.Equal(t, model.CourierStatusAvailable.String(), resp.Status)
	require.Equal(t, model.OnFoot.String(), resp.Transport)
}

func TestGetById_InvalidId_NotNumber(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/courier/abc", nil)
	req = withIDParam(req, "abc")

	rec := httptest.NewRecorder()

	h.GetById(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, ErrInvalidId.Error(), resp["error"])
}

func TestGetById_InvalidId_NonPositive(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	for _, idStr := range []string{"0", "-1"} {
		req := httptest.NewRequest(http.MethodGet, "/courier/"+idStr, nil)
		req = withIDParam(req, idStr)

		rec := httptest.NewRecorder()

		h.GetById(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, ErrInvalidId.Error(), resp["error"])
	}
}

func TestGetById_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	svc.EXPECT().
		GetCourierById(gomock.Any(), int64(1)).
		Return(nil, courier_service.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/courier/1", nil)
	req = withIDParam(req, "1")

	rec := httptest.NewRecorder()

	h.GetById(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, courier_service.ErrNotFound.Error(), resp["error"])
}

func TestGetById_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	dbErr := errors.New("db error")

	svc.EXPECT().
		GetCourierById(gomock.Any(), int64(1)).
		Return(nil, dbErr)

	req := httptest.NewRequest(http.MethodGet, "/courier/1", nil)
	req = withIDParam(req, "1")

	rec := httptest.NewRecorder()

	h.GetById(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCreateCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	reqDTO := createCourierDTO{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    string(model.CourierStatusAvailable),
		Transport: string(model.OnFoot),
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(&model.Courier{Id: 1}, nil)

	req := httptest.NewRequest(http.MethodPost, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateCourier(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]int64
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, int64(1), resp["id"])
}

func TestCreateCourier_EmptyName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	reqDTO := createCourierDTO{
		Name:      "",
		Phone:     "+79119568101",
		Status:    string(model.CourierStatusAvailable),
		Transport: string(model.OnFoot),
	}

	body, _ := json.Marshal(reqDTO)

	req := httptest.NewRequest(http.MethodPost, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, ErrEmptyName.Error(), resp["error"])
}

func TestCreateCourier_InvalidStatus(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	reqDTO := createCourierDTO{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    "invalid",
		Transport: string(model.OnFoot),
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(nil, courier_service.ErrInvalidStatus)

	req := httptest.NewRequest(http.MethodPost, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, courier_service.ErrInvalidStatus.Error(), resp["error"])
}

func TestCreateCourier_InvalidPhone(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	reqDTO := createCourierDTO{
		Name:      "Artem",
		Phone:     "123",
		Status:    string(model.CourierStatusAvailable),
		Transport: string(model.OnFoot),
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(nil, courier_service.ErrInvalidPhoneNumber)

	req := httptest.NewRequest(http.MethodPost, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, courier_service.ErrInvalidPhoneNumber.Error(), resp["error"])
}

func TestCreateCourier_InvalidTransport(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	reqDTO := createCourierDTO{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    string(model.CourierStatusAvailable),
		Transport: "invalid",
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(nil, courier_service.ErrInvalidTransport)

	req := httptest.NewRequest(http.MethodPost, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, courier_service.ErrInvalidTransport.Error(), resp["error"])
}

func TestCreateCourier_DuplicatePhone(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	reqDTO := createCourierDTO{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    string(model.CourierStatusAvailable),
		Transport: string(model.OnFoot),
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(nil, courier_service.ErrDuplicatePhone)

	req := httptest.NewRequest(http.MethodPost, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateCourier(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, courier_service.ErrDuplicatePhone.Error(), resp["error"])
}

func TestCreateCourier_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	reqDTO := createCourierDTO{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    string(model.CourierStatusAvailable),
		Transport: string(model.OnFoot),
	}

	body, _ := json.Marshal(reqDTO)

	internalErr := errors.New("db error")

	svc.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(nil, internalErr)

	req := httptest.NewRequest(http.MethodPost, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateCourier(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetAll_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	couriers := []model.Courier{
		{
			Id:        1,
			Name:      "Artem",
			Phone:     "+79119568101",
			Status:    model.CourierStatusAvailable,
			Transport: model.OnFoot,
		},
	}

	svc.EXPECT().
		GetAllCouriers(gomock.Any()).
		Return(couriers, nil)

	req := httptest.NewRequest(http.MethodGet, "/couriers", nil)
	rec := httptest.NewRecorder()

	h.GetAll(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp []courierDTO
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	require.Len(t, resp, 1)
	require.Equal(t, int64(1), resp[0].ID)
	require.Equal(t, "Artem", resp[0].Name)
	require.Equal(t, "+79119568101", resp[0].Phone)
	require.Equal(t, model.CourierStatusAvailable.String(), resp[0].Status)
	require.Equal(t, model.OnFoot.String(), resp[0].Transport)
}

func TestGetAll_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	svc.EXPECT().
		GetAllCouriers(gomock.Any()).
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/couriers", nil)
	rec := httptest.NewRecorder()

	h.GetAll(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdateCourier_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	id := int64(1)
	name := "Artem"
	phone := "+79119568101"
	status := string(model.CourierStatusAvailable)
	transport := string(model.OnFoot)

	reqDTO := updateCourierDTO{
		ID:        &id,
		Name:      &name,
		Phone:     &phone,
		Status:    &status,
		Transport: &transport,
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateCourier(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateCourier_EmptyName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	id := int64(1)
	name := ""

	reqDTO := updateCourierDTO{
		ID:   &id,
		Name: &name,
	}

	body, _ := json.Marshal(reqDTO)

	req := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, ErrEmptyName.Error(), resp["error"])
}

func TestUpdateCourier_InvalidId(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	id := int64(0)

	reqDTO := updateCourierDTO{
		ID: &id,
	}

	body, _ := json.Marshal(reqDTO)

	req := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, ErrInvalidId.Error(), resp["error"])
}

func TestUpdateCourier_InvalidStatus(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	id := int64(1)
	status := "invalid"

	reqDTO := updateCourierDTO{
		ID:     &id,
		Status: &status,
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(courier_service.ErrInvalidStatus)

	req := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, courier_service.ErrInvalidStatus.Error(), resp["error"])
}

func TestUpdateCourier_InvalidPhone(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	id := int64(1)
	phone := "123"

	reqDTO := updateCourierDTO{
		ID:    &id,
		Phone: &phone,
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(courier_service.ErrInvalidPhoneNumber)

	req := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, courier_service.ErrInvalidPhoneNumber.Error(), resp["error"])
}

func TestUpdateCourier_InvalidTransport(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	id := int64(1)
	transport := "invalid"

	reqDTO := updateCourierDTO{
		ID:        &id,
		Transport: &transport,
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(courier_service.ErrInvalidTransport)

	req := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateCourier(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, courier_service.ErrInvalidTransport.Error(), resp["error"])
}

func TestUpdateCourier_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	id := int64(1)

	reqDTO := updateCourierDTO{
		ID: &id,
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(courier_service.ErrNotFound)

	req := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateCourier(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, courier_service.ErrNotFound.Error(), resp["error"])
}

func TestUpdateCourier_DuplicatePhone(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	id := int64(1)
	phone := "+79119568101"

	reqDTO := updateCourierDTO{
		ID:    &id,
		Phone: &phone,
	}

	body, _ := json.Marshal(reqDTO)

	svc.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(courier_service.ErrDuplicatePhone)

	req := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateCourier(rec, req)

	require.Equal(t, http.StatusConflict, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, courier_service.ErrDuplicatePhone.Error(), resp["error"])
}

func TestUpdateCourier_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := courier_handler.NewMockcourierService(ctrl)
	h := NewHandler(svc)

	id := int64(1)

	reqDTO := updateCourierDTO{
		ID: &id,
	}

	body, _ := json.Marshal(reqDTO)

	internalErr := errors.New("db error")

	svc.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(internalErr)

	req := httptest.NewRequest(http.MethodPut, "/courier", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateCourier(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
