package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	nozama "nozama/src/shared"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var nozama_payments_sqs_queue = "SQSPayments"

const nozama_orders_dynamodb_table_name = "orders"

func main() {
	lambda.Start(handleNewOrder)
}

func handleNewOrder(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if req.HTTPMethod != "POST" {
		return clientHttpStatusResponse(http.StatusBadRequest, "An error occured trying to place the order")
	}

	orderRequest, err := toOrderRequest(req.Body)

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

func createOrder(newOrderRequest nozama.CreateOrderRequest) (nozama.CreateOrderEvent, error) {

	var orderItem nozama.OrderItem

	orderItem.Item = newOrderRequest.Item
	orderItem.Quantity = newOrderRequest.Quantity
	orderItem.TotalPrice = newOrderRequest.TotalPrice
	orderItem.UserID = newOrderRequest.UserID
	orderItem.Status = nozama.OrderStatusIncomplete
	orderItem.OrderID = nozama.GeneratePrimaryKey()

	err := nozama.PutItem(orderItem, nozama_orders_dynamodb_table_name)

	if err != nil {
		log.Fatalf("An error ocurred while placing order %s", err)
		return nozama.CreateOrderEvent{}, err
	}

	var newOrderEvent nozama.CreateOrderEvent
	newOrderEvent.OrderID = orderItem.OrderID
	newOrderEvent.TotalPrice = orderItem.TotalPrice

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

func sendPaymentEvent(orderEvent nozama.CreateOrderEvent) error {
	err := nozama.SendMessage(orderEvent, nozama_payments_sqs_queue)
	if err != nil {
		log.Printf("OnOrderCreatedException could not send payment event %s", err)
	}
	return err
}

func toOrderRequest(body string) (nozama.CreateOrderRequest, error) {
	b := []byte(body)
	var orderEvent nozama.CreateOrderRequest
	err := json.Unmarshal(b, &orderEvent)
	return orderEvent, err
}
