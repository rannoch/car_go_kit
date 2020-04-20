package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	service "github.com/rannoch/car/internal/service"
)

type Endpoints struct {
	GetCarEndpoint    endpoint.Endpoint
	PostCarEndpoint   endpoint.Endpoint
	PutCarEndpoint    endpoint.Endpoint
	DeleteCarEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		GetCarEndpoint:    MakeGetCarEndpoint(s),
		PostCarEndpoint:   MakePostCarEndpoint(s),
		PutCarEndpoint:    MakePutCarEndpoint(s),
		DeleteCarEndpoint: MakeDelCarEndpoint(s),
	}
}

func MakeGetCarEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetCarRequest)

		car, e := s.GetCar(ctx, req.Id)

		return GetCarResponse{
			Car: car,
			Err: e,
		}, nil
	}
}

func MakePostCarEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostCarRequest)

		e := s.PostCar(ctx, &req.Car)

		return PostCarResponse{Err: e}, nil
	}
}

func MakePutCarEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PutCarRequest)

		e := s.PutCar(ctx, req.Id, &req.Car)

		return PutCarResponse{Err: e}, nil
	}
}

func MakeDelCarEndpoint(s service.Service) endpoint.Endpoint {
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
	Car *service.Car `json:"car,omitempty"`
	Err error        `json:"err,omitempty"`
}

type PostCarRequest struct {
	Car service.Car
}

type PostCarResponse struct {
	Err error `json:"err,omitempty"`
}

type PutCarRequest struct {
	Id  int
	Car service.Car
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
