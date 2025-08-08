package main

import (
	"context"

	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/example/jsoninputguard/internal/predict"
)

var adapter *chiadapter.ChiLambda

func init() {
	adapter = chiadapter.New(predict.Router())
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return adapter.ProxyWithContext(ctx, req)
}
