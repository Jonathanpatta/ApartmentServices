package Consumers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type ConsumerHttpService struct {
	service *ConsumerService
}

func NewConsumerHttpService() (*ConsumerHttpService, error) {
	service, err := NewConsumerService()
	if err != nil {
		return nil, err
	}

	return &ConsumerHttpService{
		service: service,
	}, nil
}

func (s *ConsumerHttpService) Create(w http.ResponseWriter, r *http.Request) {
	var data Consumer
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	consumer, err := s.service.Create(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	outData, err := json.Marshal(consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(outData))
}

func (s *ConsumerHttpService) Read(w http.ResponseWriter, r *http.Request) {

	var data string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	consumer, err := s.service.Read(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	outData, err := json.Marshal(consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(outData))

}

func (s *ConsumerHttpService) Update(w http.ResponseWriter, r *http.Request) {

	var data Consumer
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	consumer, err := s.service.Update(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	outData, err := json.Marshal(consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(outData))
}

func (s *ConsumerHttpService) Delete(w http.ResponseWriter, r *http.Request) {
	var data string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	consumer, err := s.service.Delete(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	outData, err := json.Marshal(consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(outData))
}

func AddSubrouter(r *mux.Router) {
	server, err := NewConsumerHttpService()
	if err != nil {
		log.Fatal(err)
	}
	router := r.PathPrefix("/consumer").Subrouter()

	router.HandleFunc("/create", server.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/", server.Read).Methods("GET", "OPTIONS")
	router.HandleFunc("/update", server.Update).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete", server.Delete).Methods("POST", "OPTIONS")
}
