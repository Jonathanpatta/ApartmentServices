package Items

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"log"
	"net/http"
)

type ItemHttpService struct {
	service *ItemService
}

func NewItemHttpService(settings *Settings.Settings) (*ItemHttpService, error) {
	service, err := NewItemService(settings)
	if err != nil {
		return nil, err
	}

	return &ItemHttpService{
		service: service,
	}, nil
}

func (s *ItemHttpService) Create(w http.ResponseWriter, r *http.Request) {
	var data Item
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serviceId := mux.Vars(r)["serviceId"]

	item, err := s.service.Create(serviceId, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ItemHttpService) Read(w http.ResponseWriter, r *http.Request) {

	var data string
	vars := mux.Vars(r)
	data = vars["itemId"]

	item, err := s.service.Read(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ItemHttpService) Update(w http.ResponseWriter, r *http.Request) {

	var data Item
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := s.service.Update(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ItemHttpService) Delete(w http.ResponseWriter, r *http.Request) {
	var data string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := s.service.Delete(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, string(outData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ItemHttpService) List(w http.ResponseWriter, r *http.Request) {

	item, err := s.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outData, err := json.Marshal(item)
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
	server, err := NewItemHttpService(settings)
	if err != nil {
		log.Fatal(err)
	}
	router := r.PathPrefix("/item").Subrouter()

	router.HandleFunc("/list", server.List).Methods("GET", "OPTIONS")
	router.HandleFunc("/create/{serviceId}", server.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/update", server.Update).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete", server.Delete).Methods("POST", "OPTIONS")
	router.HandleFunc("/{itemId}", server.Read).Methods("GET", "OPTIONS")
}
