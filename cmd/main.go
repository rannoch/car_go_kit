package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rannoch/car"
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

	var carService car.ServiceImpl
	var db *sqlx.DB
	db, err := sqlx.Open("mysql", "root:12345678@/test?parseTime=true")
	if err != nil {
		logger.Log(err)
		return
	}

	carService.CarsRepository = car.NewRepositoryMysql(db)

	var handler http.Handler

	handler = car.MakeHTTPHandler(&carService, logger)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	logger.Log("exit", <-errs)
}