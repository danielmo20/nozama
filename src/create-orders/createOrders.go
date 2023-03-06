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

func main() {
	lambda.Start(handleCreateOrder)
}

func handleCreateOrder(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var httpErrorMessage = nozama.HttpMessageResponse{Message: "An error occurred trying to placing the order"}

	if req.HTTPMethod != http.MethodPost {
		return nozama.HttpResponse(http.StatusBadRequest, httpErrorMessage)
	}

	orderRequest, err := ToOrderRequest(req.Body)

	if err != nil {
		log.Printf("handleCreateOrder: Cannot parse order request req %s", req.Body)
		return nozama.HttpResponse(http.StatusBadRequest, httpErrorMessage)
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

	err := nozama.PutItem(orderItem, nozama.OrdersDynamoDBTableName)

	if err != nil {
		log.Printf("CreateOrder: An error ocurred while placing order. Error %s", err)
		return nozama.CreateOrderEvent{}, err
	}

	var newOrderEvent nozama.CreateOrderEvent
	newOrderEvent.OrderID = orderItem.OrderID
	newOrderEvent.TotalPrice = orderItem.TotalPrice

	return newOrderEvent, nil

}

func sendPaymentEvent(orderEvent nozama.CreateOrderEvent) error {
	err := nozama.SendMessage(orderEvent, nozama.PaymentsSQSQueue)
	if err != nil {
		log.Printf("sendPaymentEvent: could not send payment event %s", err)
	}
	return err
}

func ToOrderRequest(body string) (nozama.CreateOrderRequest, error) {
	b := []byte(body)
	var orderEvent nozama.CreateOrderRequest
	err := json.Unmarshal(b, &orderEvent)
	return orderEvent, err
}
