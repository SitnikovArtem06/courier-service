package courier_service

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/repository/courier_repository"
	"errors"
)

type CourierService struct {
	courierRepo courier_repository.CourierRepository
}

func NewCourierService(repo courier_repository.CourierRepository) *CourierService {
	return &CourierService{courierRepo: repo}
}

func (s *CourierService) CreateCourier(ctx context.Context, c *model.CreateCourierRequest) (*model.Courier, error) {

	if err := validateCreate(c); err != nil {
		return nil, err
	}

	if c.Transport == "" {
		c.Transport = model.OnFoot
	}

	courierDb := &model.CourierDB{
		Id:        0,
		Name:      c.Name,
		Phone:     c.Phone,
		Status:    c.Status,
		Transport: c.Transport,
	}

	courierDb, err := s.courierRepo.Create(ctx, courierDb)

	if err != nil {
		if errors.Is(err, courier_repository.ErrDuplicatePhoneRepo) {
			return nil, ErrDuplicatePhone
		}
		return nil, err
	}

	courier := &model.Courier{
		Id:        courierDb.Id,
		Name:      courierDb.Name,
		Phone:     courierDb.Phone,
		Status:    courierDb.Status,
		CreatedAt: courierDb.CreatedAt,
		UpdatedAt: courierDb.UpdatedAt,
		Transport: courierDb.Transport,
	}
	return courier, nil

}

func (s *CourierService) GetCourierById(ctx context.Context, id int64) (*model.Courier, error) {

	c, err := s.courierRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, courier_repository.ErrNotFoundRepo) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	resp := &model.Courier{
		Id:        c.Id,
		Name:      c.Name,
		Phone:     c.Phone,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Transport: c.Transport,
	}
	return resp, nil
}

func (s *CourierService) GetAllCouriers(ctx context.Context) ([]model.Courier, error) {

	couriersDb, err := s.courierRepo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	var couriers []model.Courier
	for _, c := range couriersDb {
		couriers = append(couriers, model.Courier{
			Id:        c.Id,
			Name:      c.Name,
			Phone:     c.Phone,
			Status:    c.Status,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Transport: c.Transport,
		})
	}

	return couriers, nil

}

func (s *CourierService) UpdateCourier(ctx context.Context, req *model.UpdateCourierRequest) error {

	if err := validateUpdate(req); err != nil {
		return err
	}

	err := s.courierRepo.Update(ctx, req)
	if errors.Is(err, courier_repository.ErrNotFoundRepo) {
		return ErrNotFound
	}
	if errors.Is(err, courier_repository.ErrDuplicatePhoneRepo) {
		return ErrDuplicatePhone
	}
	return err

}

func validateUpdate(req *model.UpdateCourierRequest) error {

	if req.Phone != nil && !validNumber(*req.Phone) {
		return ErrInvalidPhoneNumber
	}
	if req.Status != nil && !(*req.Status).IsValid() {
		return ErrInvalidStatus
	}

	if req.Transport != nil && *req.Transport != "" && !req.Transport.IsValid() {
		return ErrInvalidTransport
	}

	return nil
}

func validNumber(raw string) bool {
	buf := make([]byte, 0, len(raw))
	leadingPlus := false
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if c == '+' && len(buf) == 0 {
			leadingPlus = true
			continue
		}
		if c >= '0' && c <= '9' {
			buf = append(buf, c)
		}
	}

	if len(buf) == 11 && leadingPlus {

		return true
	}

	return false
}

func validateCreate(c *model.CreateCourierRequest) error {

	if !validNumber(c.Phone) {
		return ErrInvalidPhoneNumber
	}
	if !c.Status.IsValid() {
		return ErrInvalidStatus
	}

	if c.Transport != "" && !c.Transport.IsValid() {
		return ErrInvalidTransport
	}
	return nil
}
