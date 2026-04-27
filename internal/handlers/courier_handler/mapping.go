package courier_handler

import (
	"course-go-avito-SitnikovArtem06/internal/model"
)

func toDTO(c *model.Courier) courierDTO {
	return courierDTO{
		ID:        c.Id,
		Name:      c.Name,
		Phone:     c.Phone,
		Status:    c.Status.String(),
		Transport: c.Transport.String(),
	}
}

func toDTOs(cs []model.Courier) []courierDTO {
	out := make([]courierDTO, 0, len(cs))
	for _, c := range cs {
		out = append(out, toDTO(&c))
	}
	return out
}

func fromCreateDTO(d createCourierDTO) model.CreateCourierRequest {
	return model.CreateCourierRequest{
		Name:      d.Name,
		Phone:     d.Phone,
		Status:    model.CourierStatus(d.Status),
		Transport: model.TransportType(d.Transport),
	}
}

func fromUpdateDTO(d updateCourierDTO) model.UpdateCourierRequest {
	return model.UpdateCourierRequest{
		Id:        d.ID,
		Name:      d.Name,
		Phone:     d.Phone,
		Status:    (*model.CourierStatus)(d.Status),
		Transport: (*model.TransportType)(d.Transport),
	}
}
