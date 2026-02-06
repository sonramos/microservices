package db

import (
	"fmt"

	"github.com/sonramos/microservices/shipping/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Shipping struct {
	gorm.Model
	OrderID      int64
	Status       string
	DeliveryDays int32
	Items        []ShippingItem
}

type ShippingItem struct {
	gorm.Model
	ProductCode string
	Quantity    int32
	ShippingID  uint
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}
	err := db.AutoMigrate(&Shipping{}, ShippingItem{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (a Adapter) Get(id string) (domain.Shipping, error) {
	var entity Shipping
	res := a.db.Preload("Items").First(&entity, id)
	var items []domain.ShippingItem
	for _, item := range entity.Items {
		items = append(items, domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}
	shipping := domain.Shipping{
		ID:           int64(entity.ID),
		OrderID:      entity.OrderID,
		Status:       entity.Status,
		DeliveryDays: entity.DeliveryDays,
		Items:        items,
	}
	return shipping, res.Error
}

func (a Adapter) Save(shipping *domain.Shipping) error {
	var items []ShippingItem
	for _, item := range shipping.Items {
		items = append(items, ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}
	model := Shipping{
		OrderID:      shipping.OrderID,
		Status:       shipping.Status,
		DeliveryDays: shipping.DeliveryDays,
		Items:        items,
	}
	res := a.db.Create(&model)
	if res.Error == nil {
		shipping.ID = int64(model.ID)
	}
	return res.Error
}
