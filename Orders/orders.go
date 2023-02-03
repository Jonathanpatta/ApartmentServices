package Orders

type Order struct {
	CreatedByConsumerId string
	ItemId              string
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
