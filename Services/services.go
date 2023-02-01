package Services

import "github.com/jonathanpatta/apartmentservices/Items"

type Service struct {
	ListedItems []*Items.Item
}

type ServiceService struct {
}

func NewServiceService() (*ServiceService, error) {
	return &ServiceService{}, nil
}

func (s *ServiceService) Create(in *Service) (*Service, error) {
	return nil, nil
}

func (s *ServiceService) Read(serviceId string) (*Service, error) {
	return nil, nil
}

func (s *ServiceService) Update(in *Service) (*Service, error) {
	return nil, nil
}

func (s *ServiceService) Delete(serviceId string) (*Service, error) {
	return nil, nil
}

func (s *ServiceService) GetListedItems(serviceId string) ([]*Items.Item, error) {
	return nil, nil
}
