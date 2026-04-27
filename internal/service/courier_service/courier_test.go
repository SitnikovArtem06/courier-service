package courier_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/repository/courier_repository"
	"course-go-avito-SitnikovArtem06/internal/service/mocks"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCreate_Default(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	req := &model.CreateCourierRequest{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		Transport: "",
	}

	create := &model.CourierDB{
		Id:        0,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.OnFoot,
	}

	createModel := &model.CourierDB{
		Id:        1,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.OnFoot,
	}

	expectedModel := &model.Courier{
		Id:        1,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.OnFoot,
	}

	repo.EXPECT().
		Create(gomock.Any(), create).
		Return(createModel, nil)

	got, err := service.CreateCourier(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, expectedModel, got)
	require.Equal(t, int64(1), got.Id)
}

func TestCreate_WithTransport(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	req := &model.CreateCourierRequest{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		Transport: "car",
	}

	create := &model.CourierDB{
		Id:        0,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.Car,
	}

	createModel := &model.CourierDB{
		Id:        1,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.Car,
	}

	expectedModel := &model.Courier{
		Id:        1,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.Car,
	}

	repo.EXPECT().
		Create(gomock.Any(), create).
		Return(createModel, nil)

	got, err := service.CreateCourier(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, expectedModel, got)
	require.Equal(t, int64(1), got.Id)

}

func TestCreate_InvalidNumber(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	req := &model.CreateCourierRequest{
		Name:      "Artem",
		Phone:     "9",
		Status:    model.CourierStatusAvailable,
		Transport: "",
	}

	got, err := service.CreateCourier(context.Background(), req)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidPhoneNumber)
	require.Nil(t, got)

}

func TestCreate_InvalidStatus(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	req := &model.CreateCourierRequest{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    "availabl",
		Transport: "",
	}

	got, err := service.CreateCourier(context.Background(), req)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidStatus)
	require.Nil(t, got)

}

func TestCreate_InvalidTransport(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	req := &model.CreateCourierRequest{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		Transport: "a",
	}

	got, err := service.CreateCourier(context.Background(), req)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidTransport)
	require.Nil(t, got)

}

func TestCreate_DuplicatePhone(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	req := &model.CreateCourierRequest{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		Transport: "car",
	}

	create := &model.CourierDB{
		Id:        0,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.Car,
	}

	repo.EXPECT().Create(gomock.Any(), create).Return(nil, courier_repository.ErrDuplicatePhoneRepo)

	got, err := service.CreateCourier(context.Background(), req)

	require.ErrorIs(t, err, ErrDuplicatePhone)
	require.Nil(t, got)

}

func TestCreate_DBError(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	req := &model.CreateCourierRequest{
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		Transport: "car",
	}

	create := &model.CourierDB{
		Id:        0,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.Car,
	}

	repo.EXPECT().Create(gomock.Any(), create).Return(nil, errors.New("db error"))

	got, err := service.CreateCourier(context.Background(), req)

	require.Nil(t, got)
	require.Error(t, err)

}

func TestGet_Success(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	getModel := &model.CourierDB{
		Id:        1,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.OnFoot,
	}

	expected := &model.Courier{
		Id:        1,
		Name:      "Artem",
		Phone:     "+79119568101",
		Status:    model.CourierStatusAvailable,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Transport: model.OnFoot,
	}

	repo.EXPECT().Get(gomock.Any(), int64(1)).Return(getModel, nil)

	got, err := service.GetCourierById(context.Background(), 1)

	require.NoError(t, err)
	require.Equal(t, expected, got)

}

func TestGet_NotFound(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	repo.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, courier_repository.ErrNotFoundRepo)

	got, err := service.GetCourierById(context.Background(), 1)

	require.Nil(t, got)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGet_DBError(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	repo.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, errors.New("db error"))

	got, err := service.GetCourierById(context.Background(), 1)

	require.Nil(t, got)
	require.Error(t, err)

}

func TestUpdate_Success(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	var id int64
	id = 1
	var name string
	name = "Artem"
	var phone string
	phone = "+79119568101"
	var status model.CourierStatus
	status = model.CourierStatusAvailable
	var transport model.TransportType
	transport = model.OnFoot

	req := &model.UpdateCourierRequest{
		Id:        &id,
		Name:      &name,
		Phone:     &phone,
		Status:    &status,
		Transport: &transport,
	}

	repo.EXPECT().Update(gomock.Any(), req).Return(nil)

	err := service.UpdateCourier(context.Background(), req)

	require.Nil(t, err)

}

func TestUpdate_InvalidNumber(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	var id int64
	id = 1
	var phone string
	phone = "+7911956810"

	req := &model.UpdateCourierRequest{
		Id:    &id,
		Phone: &phone,
	}

	err := service.UpdateCourier(context.Background(), req)

	require.ErrorIs(t, ErrInvalidPhoneNumber, err)

}

func TestUpdate_InvalidStatus(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	var id int64
	id = 1
	var status string
	status = "A"

	req := &model.UpdateCourierRequest{
		Id:     &id,
		Status: (*model.CourierStatus)(&status),
	}

	err := service.UpdateCourier(context.Background(), req)

	require.ErrorIs(t, ErrInvalidStatus, err)

}

func TestUpdate_InvalidTransport(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	var id int64
	id = 1
	var transport string
	transport = "A"

	req := &model.UpdateCourierRequest{
		Id:        &id,
		Transport: (*model.TransportType)(&transport),
	}

	err := service.UpdateCourier(context.Background(), req)

	require.ErrorIs(t, ErrInvalidTransport, err)

}

func TestUpdate_NotFound(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	var id int64
	id = 1

	req := &model.UpdateCourierRequest{
		Id: &id,
	}

	repo.EXPECT().Update(gomock.Any(), req).Return(courier_repository.ErrNotFoundRepo)

	err := service.UpdateCourier(context.Background(), req)

	require.ErrorIs(t, ErrNotFound, err)

}

func TestUpdate_DuplicatePhone(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	var id int64
	id = 1

	var phone string
	phone = "+79119568101"

	req := &model.UpdateCourierRequest{
		Id:    &id,
		Phone: &phone,
	}

	repo.EXPECT().Update(gomock.Any(), req).Return(courier_repository.ErrDuplicatePhoneRepo)

	err := service.UpdateCourier(context.Background(), req)

	require.ErrorIs(t, ErrDuplicatePhone, err)

}

func TestUpdate_DBError(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	var id int64
	id = 1

	req := &model.UpdateCourierRequest{
		Id: &id,
	}

	repo.EXPECT().Update(gomock.Any(), req).Return(errors.New("db error"))

	err := service.UpdateCourier(context.Background(), req)

	require.Error(t, err)
}

func TestGetAll_Success(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	dbList := []model.CourierDB{
		{
			Id:        1,
			Name:      "Artem",
			Phone:     "+79119568101",
			Status:    model.CourierStatusAvailable,
			Transport: model.OnFoot,
		},
		{
			Id:        2,
			Name:      "Ivan",
			Phone:     "+9119568102",
			Status:    model.CourierStatusAvailable,
			Transport: model.Car,
		},
	}

	expected := []model.Courier{
		{
			Id:        1,
			Name:      "Artem",
			Phone:     "+79119568101",
			Status:    model.CourierStatusAvailable,
			Transport: model.OnFoot,
		},
		{
			Id:        2,
			Name:      "Ivan",
			Phone:     "+9119568102",
			Status:    model.CourierStatusAvailable,
			Transport: model.Car,
		},
	}

	repo.EXPECT().
		GetAll(gomock.Any()).
		Return(dbList, nil)

	got, err := service.GetAllCouriers(context.Background())

	require.NoError(t, err)
	require.Equal(t, expected, got)

}

func TestGetAll_DBError(t *testing.T) {

	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	repo := mocks.NewMockCourierRepository(ctrl)

	service := NewCourierService(repo)

	repo.EXPECT().
		GetAll(gomock.Any()).
		Return(nil, errors.New("db error"))

	got, err := service.GetAllCouriers(context.Background())

	require.Nil(t, got)
	require.Error(t, err)

}
