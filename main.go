package main

import (
	"ApartmentServices/Consumers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	//cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	//if err != nil {
	//	log.Fatalf("unable to load SDK config, %v", err)
	//}
	//
	//dynamoDbCli := dynamodb.NewFromConfig(cfg)

	router := mux.NewRouter()
	Consumers.AddSubrouter(router)
	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
