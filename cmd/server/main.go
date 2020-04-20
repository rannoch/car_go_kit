package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rannoch/car/internal/endpoint"
	"github.com/rannoch/car/internal/service"
	"github.com/rannoch/car/internal/transport"
	"github.com/rannoch/car/pb"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var carService service.ServiceImpl
	var db *sqlx.DB
	db, err := sqlx.Open("mysql", "root:1@/car_db?parseTime=true")
	if err != nil {
		_ = logger.Log(err)
		return
	}

	carService.CarsRepository = service.NewRepositoryMysql(db)

	var httpHandler http.Handler

	endpoints := endpoint.MakeServerEndpoints(&carService)

	// http transport
	httpHandler = transport.MakeHTTPHandler(endpoints, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 10)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		_ = logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, httpHandler)
	}()

	// grpc transport
	grpcPort := "8082"

	grpcServer := grpc.NewServer()
	pb.RegisterCarServiceServer(grpcServer, transport.MakeGRPCServer(endpoints, logger))

	listen, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		_ = logger.Log(err.Error())
		return
	}
	err = grpcServer.Serve(listen)

	if err != nil {
		_ = logger.Log(err.Error())
	}

	_ = logger.Log("exit", <-errs)
}
