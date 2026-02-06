package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sonramos/microservices-proto/golang/shipping"
	"github.com/sonramos/microservices/shipping/config"
	"github.com/sonramos/microservices/shipping/internal/application/core/domain"
	"github.com/sonramos/microservices/shipping/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func (a Adapter) Create(ctx context.Context, request *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	log.Println("Creating shipping...")

	var items []domain.ShippingItem
	for _, item := range request.Items {
		items = append(items, domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}
	newShipping := domain.NewShipping(request.OrderId, items)
	result, err := a.api.CreateShipping(newShipping)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to create shipping. %v", err)).Err()
	}
	return &shipping.CreateShippingResponse{
		ShippingId:   result.ID,
		DeliveryDays: result.DeliveryDays,
	}, nil
}

type Adapter struct {
	api  ports.APIPort
	port int
	shipping.UnimplementedShippingServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a Adapter) Run() {
	var err error
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}
	grpcServer := grpc.NewServer()
	shipping.RegisterShippingServer(grpcServer, a)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port ")
	}
}
