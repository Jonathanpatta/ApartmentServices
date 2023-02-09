package Subscriptions

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"log"
	"net/http"
)

type SubscriptionHttpService struct {
	service *SubscriptionService
}

func NewSubscriptionHttpService(settings *Settings.Settings) (*SubscriptionHttpService, error) {
	service, err := NewSubscriptionService(settings)
	if err != nil {
		return nil, err
	}

	return &SubscriptionHttpService{
		service: service,
	}, nil
}

func (s *SubscriptionHttpService) Create(w http.ResponseWriter, r *http.Request) {
	var data Subscription
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	consumerId := mux.Vars(r)["consumerId"]

	subscription, err := s.service.Create(consumerId, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *SubscriptionHttpService) Read(w http.ResponseWriter, r *http.Request) {

	var data string
	vars := mux.Vars(r)
	data = vars["subscriptionId"]

	subscription, err := s.service.Read(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *SubscriptionHttpService) Update(w http.ResponseWriter, r *http.Request) {

	var data Subscription
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subscription, err := s.service.Update(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *SubscriptionHttpService) Delete(w http.ResponseWriter, r *http.Request) {
	var data string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subscription, err := s.service.Delete(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *SubscriptionHttpService) List(w http.ResponseWriter, r *http.Request) {

	subscription, err := s.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddSubrouter(r *mux.Router, settings *Settings.Settings) {
	server, err := NewSubscriptionHttpService(settings)
	if err != nil {
		log.Fatal(err)
	}
	router := r.PathPrefix("/subscription").Subrouter()

	router.HandleFunc("/list", server.List).Methods("GET", "OPTIONS")
	router.HandleFunc("/create/{consumerId}", server.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/update", server.Update).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete", server.Delete).Methods("POST", "OPTIONS")
	router.HandleFunc("/{subscriptionId}", server.Read).Methods("GET", "OPTIONS")
}
