package Subscriptions

import (
	"github.com/jonathanpatta/apartmentservices/Consumers"
	"github.com/jonathanpatta/apartmentservices/Items"
)

type Subscription struct {
	CreatedBy *Consumers.Consumer
	Item      *Items.Item
}

type SubscriptionService struct {
}

func NewSubscriptionService() (*SubscriptionService, error) {
	return &SubscriptionService{}, nil
}

func (s *SubscriptionService) Create(in *Subscription) (*Subscription, error) {
	return nil, nil
}

func (s *SubscriptionService) Read(subscriptionId string) (*Subscription, error) {
	return nil, nil
}

func (s *SubscriptionService) Update(in *Subscription) (*Subscription, error) {
	return nil, nil
}

func (s *SubscriptionService) Delete(subscriptionId string) (*Subscription, error) {
	return nil, nil
}
