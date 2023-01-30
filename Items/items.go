package Items

import (
	"ApartmentServices/Producers"
	"ApartmentServices/Services"
)

type Item struct {
	CreatedBy   *Producers.Producer
	BelongsTo   *Services.Service
	Name        string
	Description string
	ImageUrls   *[]string
}

type ItemService struct {
}

func NewItemService() (*ItemService, error) {
	return &ItemService{}, nil
}

func (s *ItemService) Create(in *Item) (*Item, error) {
	return nil, nil
}

func (s *ItemService) Read(itemId string) (*Item, error) {
	return nil, nil
}

func (s *ItemService) Update(in *Item) (*Item, error) {
	return nil, nil
}

func (s *ItemService) Delete(itemId string) (*Item, error) {
	return nil, nil
}
