package shipping_adapter

import (
"context"
"log"
"time"

grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
"github.com/sonramos/microservices-proto/golang/shipping"
"github.com/sonramos/microservices/order/internal/application/core/domain"
"google.golang.org/grpc"
"google.golang.org/grpc/codes"
"google.golang.org/grpc/credentials/insecure"
"google.golang.org/grpc/status"
)

type Adapter struct {
shipping shipping.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
var opts []grpc.DialOption
opts = append(opts,
grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
grpc_retry.WithMax(5),
grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
)))
opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
conn, err := grpc.Dial(shippingServiceUrl, opts...)
if err != nil {
return nil, err
}
client := shipping.NewShippingClient(conn)
return &Adapter{shipping: client}, nil
}

func (a *Adapter) Ship(order *domain.Order) error {
var items []*shipping.ShippingItem
for _, oi := range order.OrderItems {
items = append(items, &shipping.ShippingItem{
ProductCode: oi.ProductCode,
Quantity:    oi.Quantity,
})
}
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
_, err := a.shipping.Create(ctx, &shipping.CreateShippingRequest{
OrderId: order.ID,
Items:   items,
})
if err != nil {
if code := status.Code(err); code == codes.DeadlineExceeded {
log.Printf("Shipping deadline exceeded: %v", err)
} else {
log.Printf("Failed to create shipping: %v", err)
}
}
return err
}
