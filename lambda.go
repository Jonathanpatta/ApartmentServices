package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/jonathanpatta/apartmentservices/Router"
)

var adapter *gorillamux.GorillaMuxAdapter

func LambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	fmt.Println("req:", req.Path)
	resp, err := adapter.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&req))

	return *resp.Version1(), err

}

func main() {
	router := Router.GetMainRouter()
	adapter = gorillamux.New(router)
	lambda.Start(LambdaHandler)
}
