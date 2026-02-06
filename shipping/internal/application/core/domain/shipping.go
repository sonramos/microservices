package domain

type ShippingItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int32  `json:"quantity"`
}

type Shipping struct {
	ID           int64          `json:"id"`
	OrderID      int64          `json:"order_id"`
	Status       string         `json:"status"`
	Items        []ShippingItem `json:"items"`
	DeliveryDays int32          `json:"delivery_days"`
	CreatedAt    int64          `json:"created_at"`
}

func NewShipping(orderID int64, items []ShippingItem) Shipping {
	totalQuantity := int32(0)
	for _, item := range items {
		totalQuantity += item.Quantity
	}
	days := int32(1) + totalQuantity/5
	if days < 1 {
		days = 1
	}
	return Shipping{
		OrderID:      orderID,
		Status:       "Created",
		Items:        items,
		DeliveryDays: days,
	}
}
