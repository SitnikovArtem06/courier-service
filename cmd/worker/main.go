package main

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/gateway/order"
	"course-go-avito-SitnikovArtem06/internal/handlers/queues/order/changed"
	"course-go-avito-SitnikovArtem06/internal/logger"
	"course-go-avito-SitnikovArtem06/internal/service/order_status_factory"
	"course-go-avito-SitnikovArtem06/internal/transport"
	"course-go-avito-SitnikovArtem06/pkg/kafka"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"course-go-avito-SitnikovArtem06/internal/repository/courier_repository"
	"course-go-avito-SitnikovArtem06/internal/repository/delivery_repository"
	"course-go-avito-SitnikovArtem06/internal/service/assign_service"
	"course-go-avito-SitnikovArtem06/internal/service/order_changed_service"
	"course-go-avito-SitnikovArtem06/internal/service/transport_factory"
	"course-go-avito-SitnikovArtem06/internal/tx"
	"course-go-avito-SitnikovArtem06/pkg/database"

	"github.com/joho/godotenv"
)

const TimeOut = 5 * time.Second

func run(ctx context.Context, loger logger.Logger) error {
	dbpool, err := database.InitDb(ctx)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer dbpool.Close()

	txManager := tx.NewPgxTxManager(dbpool)

	courierRepo := courier_repository.NewCourierRepo(txManager)
	deliveryRepo := delivery_repository.NewDeliveryRepository(txManager)
	transportFactory := transport_factory.NewTransportFactory()

	assignService := assign_service.NewAssignService(txManager, deliveryRepo, courierRepo, transportFactory)

	statusFactory := order_status_factory.NewOrderStatusFactory(assignService)

	baseUrl := os.Getenv("ORDER_HTTP_BASEURL")

	httpGateway := order.NewHttpGateway(baseUrl, &http.Client{Timeout: TimeOut})

	orderChanged := order_changed_service.NewOrderChangedService(statusFactory, httpGateway)

	orderChangedHandelr := changed.NewChangedHandler(orderChanged)

	kcfg, saramaCfg, err := kafka.InitKafka()
	if err != nil {
		return err
	}

	kafkaConsumer := transport.NewKafkaConsumer(kcfg.Brokers, kcfg.Topic, orderChangedHandelr, saramaCfg)

	errCh := make(chan error, 1)

	go func() {
		if err := kafkaConsumer.Run(ctx); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():

		loger.Log("Shutting down worker")

		return nil
	case err := <-errCh:
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return fmt.Errorf("worker: %w", err)
	}
}

func main() {
	_ = godotenv.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	loger := logger.NewLogger()

	if err := run(ctx, loger); err != nil {
		loger.Log(fmt.Sprintf("fatal: %v", err))
		os.Exit(1)
	}
}
