package Consumers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"log"
	"net/http"
)

type ConsumerHttpService struct {
	service *ConsumerService
}

func NewConsumerHttpService(settings *Settings.Settings) (*ConsumerHttpService, error) {
	service, err := NewConsumerService(settings)
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
		return
	}

	consumer, err := s.service.Create(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ConsumerHttpService) Read(w http.ResponseWriter, r *http.Request) {

	var data string
	//err := json.NewDecoder(r.Body).Decode(&data)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	data = r.URL.Query().Get("id")

	consumer, err := s.service.Read(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ConsumerHttpService) Update(w http.ResponseWriter, r *http.Request) {

	var data Consumer
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	consumer, err := s.service.Update(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(consumer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ConsumerHttpService) Delete(w http.ResponseWriter, r *http.Request) {
	var data string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	consumer, err := s.service.Delete(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(consumer)
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
	server, err := NewConsumerHttpService(settings)
	if err != nil {
		log.Fatal(err)
	}
	router := r.PathPrefix("/consumer").Subrouter()

	router.HandleFunc("/create", server.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/", server.Read).Methods("GET", "OPTIONS")
	router.HandleFunc("/update", server.Update).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete", server.Delete).Methods("POST", "OPTIONS")
}
