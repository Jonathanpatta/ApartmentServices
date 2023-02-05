package Producers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"log"
	"net/http"
)

type ProducerHttpService struct {
	service *ProducerService
}

func NewProducerHttpService(settings *Settings.Settings) (*ProducerHttpService, error) {
	service, err := NewProducerService(settings)
	if err != nil {
		return nil, err
	}

	return &ProducerHttpService{
		service: service,
	}, nil
}

func (s *ProducerHttpService) Create(w http.ResponseWriter, r *http.Request) {
	var data Producer
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	producer, err := s.service.Create(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(producer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ProducerHttpService) Read(w http.ResponseWriter, r *http.Request) {

	var data string
	vars := mux.Vars(r)
	data = vars["producerId"]

	producer, err := s.service.Read(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(producer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ProducerHttpService) Update(w http.ResponseWriter, r *http.Request) {

	var data Producer
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	producer, err := s.service.Update(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(producer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ProducerHttpService) Delete(w http.ResponseWriter, r *http.Request) {
	var data string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	producer, err := s.service.Delete(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(producer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ProducerHttpService) List(w http.ResponseWriter, r *http.Request) {

	producer, err := s.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(producer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ProducerHttpService) GetServices(w http.ResponseWriter, r *http.Request) {

	producerId := mux.Vars(r)["producerId"]

	producer, err := s.service.GetServices(producerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(producer)
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
	server, err := NewProducerHttpService(settings)
	if err != nil {
		log.Fatal(err)
	}
	router := r.PathPrefix("/producer").Subrouter()

	router.HandleFunc("/list", server.List).Methods("GET", "OPTIONS")
	router.HandleFunc("/create", server.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/update", server.Update).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete", server.Delete).Methods("POST", "OPTIONS")
	router.HandleFunc("/{producerId}/services", server.GetServices).Methods("GET", "OPTIONS")
	router.HandleFunc("/{producerId}", server.Read).Methods("GET", "OPTIONS")
}
