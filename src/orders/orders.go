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

var nozama_payments_sqs_queue = "SQSPayments"

const nozama_orders_dynamodb_table_name = "orders"

func main() {
	lambda.Start(handleNewOrder)
}

type CreateOrderRequest struct {
	OrderID    string `json:"order_id"`
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

	if req.HTTPMethod != "POST" {
		return clientHttpStatusResponse(http.StatusBadRequest, "An error occured trying to place the order")
	}

	orderRequest, err := getOrderRequest(req.Body)

	if err != nil {
		log.Printf("NOZAMA - Cannot get order request req %s", req.Body)
	}

	orderEvent, err := createOrder(orderRequest)

	if err != nil {
		return clientHttpStatusResponse(http.StatusBadRequest, "An error occured trying to place the order")
	}

	err = sendPaymentEvent(orderEvent)

	if err != nil {
		return clientHttpStatusResponse(http.StatusBadRequest, "An error occured trying to place the order")
	} else {
		return clientHttpStatusResponse(http.StatusCreated, "Order Created")
	}

}

func createOrder(newOrder CreateOrderRequest) (CreateOrderEvent, error) {
	newOrder.OrderID = nozama.GeneratePrimaryKey()
	err := nozama.PutItem(newOrder, nozama_orders_dynamodb_table_name)

	if err != nil {
		log.Fatalf("An error ocurred while placing order %s", err)
		return CreateOrderEvent{}, err
	}

	var newOrderEvent CreateOrderEvent
	newOrderEvent.OrderID = newOrder.OrderID
	newOrderEvent.TotalPrice = newOrder.TotalPrice

	return newOrderEvent, nil

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
	err := nozama.SendMessage(orderEvent.OrderID, nozama_payments_sqs_queue)
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
