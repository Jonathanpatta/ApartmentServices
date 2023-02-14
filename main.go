package main

import (
	"fmt"
	"github.com/jonathanpatta/apartmentservices/Router"
	"net/http"
)

func main() {
	router := Router.GetMainRouter()
	http.Handle("/", router)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
