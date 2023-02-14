package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/jonathanpatta/apartmentservices/Consumers"
	"github.com/jonathanpatta/apartmentservices/Items"
	"github.com/jonathanpatta/apartmentservices/Middleware"
	"github.com/jonathanpatta/apartmentservices/Orders"
	"github.com/jonathanpatta/apartmentservices/Producers"
	"github.com/jonathanpatta/apartmentservices/Services"
	"github.com/jonathanpatta/apartmentservices/Settings"
	"github.com/jonathanpatta/apartmentservices/Subscriptions"
	"log"
)

var adapter *gorillamux.GorillaMuxAdapter

func LambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	fmt.Println("req:", req.Path)
	resp, err := adapter.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&req))

	return *resp.Version1(), err

}

func main() {

	settings, err := Settings.NewSettings()
	if err != nil {
		log.Fatalf("unable to load settings, %v", err)
	}

	router := mux.NewRouter()
	Consumers.AddSubrouter(router, settings)
	Producers.AddSubrouter(router, settings)
	Services.AddSubrouter(router, settings)
	Items.AddSubrouter(router, settings)
	Orders.AddSubrouter(router, settings)
	Subscriptions.AddSubrouter(router, settings)
	router.Use(Middleware.CorsMiddleware)

	//http.Handle("/", router)
	//http.ListenAndServe(":8000", nil)

	adapter = gorillamux.New(router)
	lambda.Start(LambdaHandler)
}
