package Producers

import "github.com/jonathanpatta/apartmentservices/Services"

type Producer struct {
	Services []*Services.Service
}

type ProducerService struct {
}

func NewProducerService() (*ProducerService, error) {
	return &ProducerService{}, nil
}

func (s *ProducerService) Create(in *Producer) (*Producer, error) {
	return nil, nil
}

func (s *ProducerService) Read(producerId string) (*Producer, error) {
	return nil, nil
}

func (s *ProducerService) Update(in *Producer) (*Producer, error) {
	return nil, nil
}

func (s *ProducerService) Delete(producerId string) (*Producer, error) {
	return nil, nil
}

func (s *ProducerService) GetServices(producerId string) ([]*Services.Service, error) {
	return nil, nil
}
