package api

import (
"github.com/sonramos/microservices/order/internal/application/core/domain"
"github.com/sonramos/microservices/order/internal/ports"
"google.golang.org/grpc/codes"
"google.golang.org/grpc/status"
)

type Application struct {
db       ports.DBPort
payment  ports.PaymentPort
shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
return &Application{
db:       db,
payment:  payment,
shipping: shipping,
}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
// Validate total quantity of items (max 50)
var totalQuantity int32
for _, item := range order.OrderItems {
totalQuantity += item.Quantity
}
if totalQuantity > 50 {
return domain.Order{}, status.Errorf(codes.InvalidArgument, "Order cannot have more than 50 items in total (got %d).", totalQuantity)
}

// Validate stock - check all product codes exist
var productCodes []string
for _, item := range order.OrderItems {
productCodes = append(productCodes, item.ProductCode)
}
if err := a.db.ValidateStock(productCodes); err != nil {
return domain.Order{}, status.Errorf(codes.NotFound, "Stock validation failed: %v", err)
}

// Save order
err := a.db.Save(&order)
if err != nil {
return domain.Order{}, err
}

// Charge payment
paymentErr := a.payment.Charge(&order)
if paymentErr != nil {
return domain.Order{}, paymentErr
}

// Ship order (only after successful payment)
shipErr := a.shipping.Ship(&order)
if shipErr != nil {
return domain.Order{}, shipErr
}

return order, nil
}
