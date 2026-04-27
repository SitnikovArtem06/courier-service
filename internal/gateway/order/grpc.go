package order

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/model"
	"course-go-avito-SitnikovArtem06/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"time"
)

type GrpcGateway struct {
	client pb.OrdersServiceClient
}

func (g *GrpcGateway) GetNewOrders(ctx context.Context, from time.Time) (*model.OrdersResponse, error) {

	req := pb.GetOrdersRequest{From: timestamppb.New(from)}

	resp, err := g.client.GetOrders(ctx, &req)
	if err != nil {
		return nil, err
	}

	ordersId := make([]string, 0, len(resp.Orders))
	createdAt := make([]time.Time, 0, len(resp.Orders))

	for _, order := range resp.Orders {
		ordersId = append(ordersId, order.Id)
		createdAt = append(createdAt, order.CreatedAt.AsTime())

	}

	return &model.OrdersResponse{OrdersId: ordersId, CreatedAt: createdAt}, nil
}

func NewGrpcGateway() (*GrpcGateway, error) {
	addr := os.Getenv("ORDER_GRPC_HOST")

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewOrdersServiceClient(conn)

	return &GrpcGateway{
		client: client,
	}, nil
}
