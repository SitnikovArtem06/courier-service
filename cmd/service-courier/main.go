package main

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/handlers"
	"course-go-avito-SitnikovArtem06/internal/handlers/assign_handler"
	"course-go-avito-SitnikovArtem06/internal/handlers/courier_handler"
	logger "course-go-avito-SitnikovArtem06/internal/logger"
	"course-go-avito-SitnikovArtem06/internal/middleware"
	"course-go-avito-SitnikovArtem06/internal/middleware/ratelimiter"
	"course-go-avito-SitnikovArtem06/internal/observability"
	"course-go-avito-SitnikovArtem06/internal/repository/courier_repository"
	"course-go-avito-SitnikovArtem06/internal/repository/delivery_repository"
	"course-go-avito-SitnikovArtem06/internal/service/assign_service"
	"course-go-avito-SitnikovArtem06/internal/service/courier_service"
	"course-go-avito-SitnikovArtem06/internal/service/delivery_monitor_service"
	"course-go-avito-SitnikovArtem06/internal/service/transport_factory"
	"course-go-avito-SitnikovArtem06/internal/tx"
	"course-go-avito-SitnikovArtem06/pkg/database"
	"errors"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"log"
	"net/http"
	"os"

	"github.com/spf13/pflag"
)

const TimeOut = 5
const Capacity = 5.0
const Refill = 5.0

func run(ctx context.Context, port string, timesec int, loger logger.Logger) error {

	dbpool, err := database.InitDb(ctx)

	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	defer dbpool.Close()

	txManager := tx.NewPgxTxManager(dbpool)

	repo := courier_repository.NewCourierRepo(txManager)
	courierService := courier_service.NewCourierService(repo)

	transportFactory := transport_factory.NewTransportFactory()

	deliveryRepo := delivery_repository.NewDeliveryRepository(txManager)
	assignService := assign_service.NewAssignService(txManager, deliveryRepo, repo, transportFactory)

	assignHandler := assign_handler.NewAssignHandler(assignService)

	handler := courier_handler.NewHandler(courierService)

	interval := time.Duration(timesec) * time.Second

	monitorService := delivery_monitor_service.NewDeliveryMonitorService(deliveryRepo, repo, interval)

	// gateway, err := order.NewGrpcGateway()
	//if err != nil {
	//return fmt.Errorf("fail with grpc %w", err)
	//}

	// monitorOrder := order_monitor_service.NewOrderMonitorService(gateway, assignService, 5*time.Second)

	errMonitorCh := make(chan error, 1)

	go func() {
		if err := monitorService.MonitorDeadline(ctx); err != nil {
			errMonitorCh <- err
		}
	}()

	errMonitorOrderCh := make(chan error, 1)

	// go func() {
	//if err := monitorOrder.Monitor(ctx); err != nil {
	//errMonitorOrderCh <- err
	//}
	//}()

	observability.Register()

	r := handlers.Routes(handler, assignHandler)

	tokenBucket := ratelimiter.NewTokenBucket(Capacity, Refill)

	rLimiter := ratelimiter.RateLimiterMiddleware(tokenBucket, loger)

	rMiddleware := middleware.ObservabilityMiddleware(rLimiter(r), loger)

	pprofSrv := observability.StartPprof("0.0.0.0:6060", loger)
	defer observability.StopServer(pprofSrv)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: rMiddleware,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("Server start on : %s", port)
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():

		shutdownCtx, cancel := context.WithTimeout(context.Background(), TimeOut*time.Second)

		defer cancel()

		log.Println("Shutting down service-courier")

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
		return nil

	case err = <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("listen: %w", err)
	case err = <-errMonitorCh:
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return fmt.Errorf("monitor delivery: %w", err)
	case err = <-errMonitorOrderCh:
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return fmt.Errorf("monitor orders: %w", err)

	}
}

func main() {

	log.SetOutput(os.Stdout)

	_ = godotenv.Load()

	port := os.Getenv("PORT")
	monitorTime := os.Getenv("MONITOR_TIME")

	timesec, _ := strconv.Atoi(monitorTime)

	var portFlag = pflag.String("port", port, "Server port")

	pflag.Parse()

	port = *portFlag

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	loger := logger.NewLogger()

	if err := run(ctx, port, timesec, loger); err != nil {
		loger.Log(fmt.Sprintf("fatal: %v", err))
		os.Exit(1)
	}
}
