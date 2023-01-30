package Orders

import (
	"ApartmentServices/Consumers"
	"ApartmentServices/Items"
)

type Order struct {
	CreatedBy *Consumers.Consumer
	Item      *Items.Item
}

type OrderService struct {
}

func NewOrderService() (*OrderService, error) {
	return &OrderService{}, nil
}

func (s *OrderService) Create(in *Order) (*Order, error) {
	return nil, nil
}

func (s *OrderService) Read(orderId string) (*Order, error) {
	return nil, nil
}

func (s *OrderService) Update(in *Order) (*Order, error) {
	return nil, nil
}

func (s *OrderService) Delete(orderId string) (*Order, error) {
	return nil, nil
}
