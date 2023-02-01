package Consumers

import "fmt"

type Consumer struct {
}

type ConsumerService struct {
}

func NewConsumerService() (*ConsumerService, error) {
	return &ConsumerService{}, nil
}

func (s *ConsumerService) Create(in *Consumer) (*Consumer, error) {
	return nil, nil
}

func (s *ConsumerService) Read(consumerId string) (*Consumer, error) {
	fmt.Println("hi from read consumer")
	return nil, nil
}

func (s *ConsumerService) Update(in *Consumer) (*Consumer, error) {
	return nil, nil
}

func (s *ConsumerService) Delete(consumerId string) (*Consumer, error) {
	return nil, nil
}
