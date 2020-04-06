package car

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetCarEndpoint    endpoint.Endpoint
	PostCarEndpoint   endpoint.Endpoint
	PutCarEndpoint    endpoint.Endpoint
	DeleteCarEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetCarEndpoint:    MakeGetCarEndpoint(s),
		PostCarEndpoint:   MakePostCarEndpoint(s),
		PutCarEndpoint:    MakePutCarEndpoint(s),
		DeleteCarEndpoint: MakeDelCarEndpoint(s),
	}
}

func MakeGetCarEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetCarRequest)

		car, e := s.GetCar(ctx, req.Id)

		return GetCarResponse{
			Car: car,
			Err: e,
		}, nil
	}
}

func MakePostCarEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostCarRequest)

		e := s.PostCar(ctx, &req.Car)

		return PostCarResponse{Err: e}, nil
	}
}

func MakePutCarEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PutCarRequest)

		e := s.PutCar(ctx, req.Id, &req.Car)

		return PutCarResponse{Err: e}, nil
	}
}

func MakeDelCarEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DelCarRequest)

		e := s.DeleteCar(ctx, req.Id)

		return DelCarResponse{Err: e}, nil
	}
}

type GetCarRequest struct {
	Id int
}

type GetCarResponse struct {
	Car *Car  `json:"car,omitempty"`
	Err error `json:"err,omitempty"`
}

type PostCarRequest struct {
	Car Car
}

type PostCarResponse struct {
	Err error `json:"err,omitempty"`
}

type PutCarRequest struct {
	Id  int
	Car Car
}

type PutCarResponse struct {
	Err error `json:"err,omitempty"`
}

type DelCarRequest struct {
	Id int
}

type DelCarResponse struct {
	Err error `json:"err,omitempty"`
}
