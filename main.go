package main

import (
	"context"
	"fmt"

	"microservice_gokit_base/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"microservice_gokit_base/src/application/endpoints"

	appHttp "microservice_gokit_base/src/application/transport/http"
	"microservice_gokit_base/src/domain/model"
	domainRepo "microservice_gokit_base/src/domain/repository"
	domainSvc "microservice_gokit_base/src/domain/service"
	"microservice_gokit_base/src/domain/utils"
	infraRepo "microservice_gokit_base/src/infraestructure/repository"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {

	//Set UTC for time management
	os.Setenv("TZ", "Europe/Spain")
	var (
		config     = config.Instance()
		ctx        = context.Background()
		apiVersion = "/api/v1/"
		httpAddr   = ":" + config.HTTPPort
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"name", "ms-base",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	var uuidGen utils.IUUIDGenerator
	{
		uuidGen = utils.NewUUIDGenerator()
	}

	var dateGen utils.IDateGenerator
	{
		dateGen = utils.NewDateGenerator()
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var repo domainRepo.IOrderRepository
	{
		if config.DB == "mongo" {
			connection := infraRepo.GetConnectionMongo(ctx, logger)
			r, err := infraRepo.NewOrderMongoRepository(connection, logger)
			if err != nil {
				level.Error(logger).Log("exit", err)
				os.Exit(-1)
			}
			repo = r
		} else {
			var dbMemory = &infraRepo.DbMemory{
				Data: []model.Order{},
			}
			r, err := infraRepo.NewOrderRepositoryMem(dbMemory, logger)
			if err != nil {
				level.Error(logger).Log("exit", err)
				os.Exit(-1)
			}
			repo = r
		}
	}

	// Create Order Services
	var svc domainSvc.IOrderService
	{
		svc = domainSvc.NewOrderService(repo, uuidGen, dateGen, logger)
	}

	var orderHandler http.Handler
	{
		endpoints := endpoints.MakeOrderEndpoints(svc)
		orderHandler = appHttp.NewHTTPOrder(endpoints, logger, apiVersion)
	}

	mux := http.NewServeMux()
	mux.Handle(apiVersion, orderHandler)
	http.Handle("/", accessControl(mux))

	// INIT WEB SERVER
	errs := make(chan error, 2)

	go func() {
		level.Info(logger).Log("transport", "http", "address", httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(httpAddr, nil)
	}()

	// HANDLE OS FINISH SIGNAL
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	level.Error(logger).Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
