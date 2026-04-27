package transport_factory

import (
	"course-go-avito-SitnikovArtem06/internal/model"
	"time"
)

type Transport interface {
	Deadline() time.Time
}
type OnFoot struct{}

func (OnFoot) Deadline() time.Time {
	return time.Now().UTC().Add(30 * time.Minute)
}

type Scooter struct{}

func (Scooter) Deadline() time.Time {
	return time.Now().UTC().Add(15 * time.Minute)
}

type Car struct{}

func (Car) Deadline() time.Time {
	return time.Now().UTC().Add(5 * time.Minute)
}

type TransportFactory interface {
	Get(transportType model.TransportType) Transport
}

type TransportFactoryImpl struct {
}

func NewTransportFactory() *TransportFactoryImpl {
	return &TransportFactoryImpl{}
}
func (f *TransportFactoryImpl) Get(transportType model.TransportType) Transport {
	switch transportType {
	case model.OnFoot:
		return OnFoot{}
	case model.Scooter:
		return Scooter{}
	case model.Car:
		return Car{}
	default:
		return OnFoot{}
	}
}
