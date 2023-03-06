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

	var httpErrorMessage = nozama.HttpMessageResponse{Message: "An error occured trying to placte the oreder"}

	if req.HTTPMethod != http.MethodPost {
		return nozama.HttpResponse(http.StatusBadRequest, httpErrorMessage)
	}

	orderRequest, err := ToOrderRequest(req.Body)

	if err != nil {
		log.Printf("NOZAMA - Cannot get order request req %s", req.Body)
	}

	orderEvent, err := CreateOrder(orderRequest)

	if err != nil {
		return nozama.HttpResponse(http.StatusBadRequest, httpErrorMessage)
	}

	err = sendPaymentEvent(orderEvent)

	if err != nil {
		return nozama.HttpResponse(http.StatusBadRequest, httpErrorMessage)
	} else {
		return nozama.HttpResponse(http.StatusCreated, orderEvent)
	}

}

func CreateOrder(newOrderRequest nozama.CreateOrderRequest) (nozama.CreateOrderEvent, error) {

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

func sendPaymentEvent(orderEvent nozama.CreateOrderEvent) error {
	err := nozama.SendMessage(orderEvent, nozama_payments_sqs_queue)
	if err != nil {
		log.Printf("OnOrderCreatedException could not send payment event %s", err)
	}
	return err
}

func ToOrderRequest(body string) (nozama.CreateOrderRequest, error) {
	b := []byte(body)
	var orderEvent nozama.CreateOrderRequest
	err := json.Unmarshal(b, &orderEvent)
	return orderEvent, err
}
