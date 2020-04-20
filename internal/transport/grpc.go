package transport

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rannoch/car/internal/endpoint"
	"github.com/rannoch/car/internal/service"
	pb "github.com/rannoch/car/pb"
)

type grpcServer struct {
	getCar    grpctransport.Handler
	postCar   grpctransport.Handler
	putCar    grpctransport.Handler
	deleteCar grpctransport.Handler
}

func (g grpcServer) GetCar(ctx context.Context, requestPb *pb.GetCarRequestPb) (*pb.CarPb, error) {
	_, resp, err := g.getCar.ServeGRPC(ctx, requestPb)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.CarPb), nil
}

func (g grpcServer) PostCar(ctx context.Context, idPb *pb.CarWithoutIdPb) (*empty.Empty, error) {
	_, resp, err := g.postCar.ServeGRPC(ctx, idPb)
	if err != nil {
		return nil, err
	}
	return resp.(*empty.Empty), nil
}

func (g grpcServer) PutCar(ctx context.Context, carPb *pb.CarPb) (*empty.Empty, error) {
	_, resp, err := g.putCar.ServeGRPC(ctx, carPb)
	if err != nil {
		return nil, err
	}
	return resp.(*empty.Empty), nil
}

func (g grpcServer) DeleteCar(ctx context.Context, requestPb *pb.DeleteCarRequestPb) (*empty.Empty, error) {
	_, resp, err := g.deleteCar.ServeGRPC(ctx, requestPb)
	if err != nil {
		return nil, err
	}
	return resp.(*empty.Empty), nil

}

func MakeGRPCServer(endpoints endpoint.Endpoints, logger log.Logger) pb.CarServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		getCar: grpctransport.NewServer(
			endpoints.GetCarEndpoint,
			decodeGRPCGetCarRequest,
			encodeGRPCGetCarRequest,
			options...,
		),
		postCar: grpctransport.NewServer(
			endpoints.PostCarEndpoint,
			decodeGRPCPostCarRequest,
			encodeGRPCPostCarResponse,
			options...,
		),
		putCar: grpctransport.NewServer(
			endpoints.PutCarEndpoint,
			decodeGRPCPutCarRequest,
			encodeGRPCPutCarResponse,
			options...,
		),
		deleteCar: grpctransport.NewServer(
			endpoints.DeleteCarEndpoint,
			decodeGRPCDelCarRequest,
			encodeGRPCDelCarResponse,
			options...,
		),
	}
}

func decodeGRPCGetCarRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetCarRequestPb)
	return endpoint.GetCarRequest{
		Id: int(req.GetId()),
	}, nil
}

func encodeGRPCGetCarRequest(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.GetCarResponse)

	if resp.Err != nil {
		return nil, resp.Err
	}

	carCreatedTimestampProto, err := ptypes.TimestampProto(resp.Car.Created)
	if err != nil {
		return nil, err
	}

	return &pb.CarPb{
		Id:      int64(resp.Car.Id),
		Brand:   resp.Car.Brand,
		Model:   resp.Car.Model,
		Created: carCreatedTimestampProto,
	}, nil
}

func decodeGRPCPostCarRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CarWithoutIdPb)

	created, err := ptypes.Timestamp(req.Created)
	if err != nil {
		return &empty.Empty{}, err
	}

	return endpoint.PostCarRequest{
		Car: service.Car{
			Brand:   req.Brand,
			Model:   req.Model,
			Created: created,
		},
	}, nil
}

func encodeGRPCPostCarResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.PostCarResponse)

	return &empty.Empty{}, resp.Err
}

func decodeGRPCPutCarRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CarPb)

	created, err := ptypes.Timestamp(req.Created)
	if err != nil {
		return &empty.Empty{}, err
	}

	return endpoint.PutCarRequest{
		Id: int(req.Id),
		Car: service.Car{
			Id:      int(req.Id),
			Brand:   req.Brand,
			Model:   req.Model,
			Created: created,
		},
	}, nil
}

func encodeGRPCPutCarResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.PutCarResponse)

	return &empty.Empty{}, resp.Err
}

func decodeGRPCDelCarRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteCarRequestPb)

	return endpoint.DelCarRequest{Id: int(req.Id)}, nil
}

func encodeGRPCDelCarResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.DelCarResponse)

	if resp.Err != nil {
		return nil, resp.Err
	}

	return &empty.Empty{}, nil
}
