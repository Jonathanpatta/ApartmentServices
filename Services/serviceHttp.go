package Services

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"log"
	"net/http"
)

type ServiceHttpService struct {
	service *ServiceService
}

func NewServiceHttpService(settings *Settings.Settings) (*ServiceHttpService, error) {
	service, err := NewServiceService(settings)
	if err != nil {
		return nil, err
	}

	return &ServiceHttpService{
		service: service,
	}, nil
}

func (s *ServiceHttpService) Create(w http.ResponseWriter, r *http.Request) {
	var data Service
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	producerId := mux.Vars(r)["producerId"]

	service, err := s.service.Create(producerId, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ServiceHttpService) Read(w http.ResponseWriter, r *http.Request) {

	var data string
	vars := mux.Vars(r)
	data = vars["serviceId"]

	service, err := s.service.Read(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ServiceHttpService) Update(w http.ResponseWriter, r *http.Request) {

	var data Service
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service, err := s.service.Update(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ServiceHttpService) Delete(w http.ResponseWriter, r *http.Request) {
	var data string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service, err := s.service.Delete(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ServiceHttpService) List(w http.ResponseWriter, r *http.Request) {

	service, err := s.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ServiceHttpService) GetItems(w http.ResponseWriter, r *http.Request) {

	serviceId := mux.Vars(r)["serviceId"]

	producer, err := s.service.GetItems(serviceId)
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
	server, err := NewServiceHttpService(settings)
	if err != nil {
		log.Fatal(err)
	}
	router := r.PathPrefix("/service").Subrouter()

	router.HandleFunc("/list", server.List).Methods("GET", "OPTIONS")
	router.HandleFunc("/create/{producerId}", server.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/update", server.Update).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete", server.Delete).Methods("POST", "OPTIONS")
	router.HandleFunc("/{serviceId}", server.Read).Methods("GET", "OPTIONS")
	router.HandleFunc("/{serviceId}/items", server.GetItems).Methods("GET", "OPTIONS")

}
