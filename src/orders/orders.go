package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	nozama "nozama/src"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handleNewOrder)
}

type CreateOrderRequest struct {
	UserID     string `json:"user_id"`
	Item       string `json:"item"`
	Quantity   int    `json:"quantity"`
	TotalPrice int64  `json:"total_price"`
}

type CreateOrderEvent struct {
	OrderID    string `json:"order_id"`
	TotalPrice int64  `json:"total_price"`
}

func handleNewOrder(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received req %#v", req)

	orderRequest, err := getOrderRequest(req.Body)

	if err != nil {
		log.Printf("NOZAMA - Cannot get order request req %s", req.Body)
	}

	var orderEvent CreateOrderEvent
	orderEvent.OrderID = "19870106"
	orderEvent.TotalPrice = orderRequest.TotalPrice

	err = sendPaymentEvent(orderEvent)

	if err != nil {
		return clientHttpStatusResponse(http.StatusBadRequest, "An error occured trying to place the order")
	} else {
		return clientHttpStatusResponse(http.StatusAccepted, "Order Created")
	}

}

func clientHttpStatusResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {

	if message == events.NewNullAttribute().String() {
		message = http.StatusText(statusCode)
	}

	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: statusCode,
	}, nil
}

func sendPaymentEvent(orderEvent CreateOrderEvent) error {
	err := nozama.SendMessage(orderEvent.OrderID)
	if err != nil {
		log.Printf("OnOrderCreatedException could not send payment event %s", err)
	}
	return err
}

func getOrderRequest(body string) (CreateOrderRequest, error) {
	b := []byte(body)
	var orderEvent CreateOrderRequest
	err := json.Unmarshal(b, &orderEvent)
	return orderEvent, err
}
