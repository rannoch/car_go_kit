package car

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/rannoch/car/internal/service"
	"github.com/rannoch/car/pb"
	"time"
)

type GrpcBinding struct {
	service.Service
}

func (g GrpcBinding) PostCar(ctx context.Context, pb *pb.CarWithoutIdPb) (*empty.Empty, error) {
	created, err := ptypes.Timestamp(pb.Created)
	if err != nil {
		return &empty.Empty{}, err
	}

	car := service.Car{
		Brand:   pb.Brand,
		Model:   pb.Model,
		Created: created,
	}

	err = g.Service.PostCar(ctx, &car)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (g GrpcBinding) PutCar(ctx context.Context, pb *pb.CarPb) (*empty.Empty, error) {
	car := service.Car{
		Brand:   pb.Brand,
		Model:   pb.Model,
		Created: time.Time{},
	}

	err := g.Service.PutCar(ctx, int(pb.Id), &car)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (g GrpcBinding) DeleteCar(ctx context.Context, pb *pb.DeleteCarRequestPb) (*empty.Empty, error) {
	err := g.Service.DeleteCar(ctx, int(pb.Id))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (g GrpcBinding) GetCar(ctx context.Context, getCarRequestPb *pb.GetCarRequestPb) (*pb.CarPb, error) {
	car, err := g.Service.GetCar(ctx, int(getCarRequestPb.Id))
	if err != nil {
		return nil, err
	}

	carCreatedTimestampProto, err := ptypes.TimestampProto(car.Created)
	if err != nil {
		return nil, err
	}

	return &pb.CarPb{
		Id:      int64(car.Id),
		Brand:   car.Brand,
		Model:   car.Model,
		Created: carCreatedTimestampProto,
	}, nil
}
