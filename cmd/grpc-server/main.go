package main

import (
	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rannoch/car"
	"github.com/rannoch/car/internal/service"
	"github.com/rannoch/car/pb"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	port := "8080"

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	listen, err := net.Listen("tcp", ":"+port)

	if err != nil {
		_ = logger.Log(err)
		return
	}

	var carService service.ServiceImpl
	var db *sqlx.DB
	db, err = sqlx.Open("mysql", "root:1@/car_db?parseTime=true")
	if err != nil {
		_ = logger.Log(err.Error())
		return
	}

	carService.CarsRepository = service.NewRepositoryMysql(db)

	server := grpc.NewServer()
	pb.RegisterCarServiceServer(server, car.GrpcBinding{Service: &carService})

	err = server.Serve(listen)

	if err != nil {
		_ = logger.Log(err.Error())
	}
}
