package courier_handler

type courierDTO struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	Transport string `json:"transport_type"`
}

type createCourierDTO struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	Transport string `json:"transport_type"`
}

type updateCourierDTO struct {
	ID        *int64  `json:"id"`
	Name      *string `json:"name,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Status    *string `json:"status,omitempty"`
	Transport *string `json:"transport_type,omitempty"`
}

func (d updateCourierDTO) validateUpdate() error {
	if d.ID == nil || *d.ID <= 0 {
		return ErrInvalidId
	}

	if d.Name != nil && *d.Name == "" {
		return ErrEmptyName
	}
	return nil
}

func (c createCourierDTO) validateCreate() error {
	if c.Name == "" {
		return ErrEmptyName
	}
	return nil
}
