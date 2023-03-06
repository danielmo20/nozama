package main

import (
	"context"
	"encoding/json"
	"log"
	nozama "nozama/src/shared"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleCreatePayment(ctx context.Context, sqsEvent events.SQSEvent) error {

	messageBody := sqsEvent.Records[0].Body

	log.Printf("handleCreatePayment: event received: %s ", messageBody)

	createOrderEvent, err := ToCreateOrderEvent(messageBody)

	if err != nil {
		log.Printf("handleCreatePayment - An error while parsing message. Error %s", err)
		return err
	}

	err = CreatePayment(createOrderEvent)

	if err != nil {
		log.Printf("handleCreatePayment: An error while creating a Payment. Error %s", err)
		return err
	}

	return err
}

func main() {
	lambda.Start(handleCreatePayment)
}

func ToCreateOrderEvent(body string) (nozama.CreateOrderEvent, error) {
	b := []byte(body)
	var orderEvent nozama.CreateOrderEvent
	err := json.Unmarshal(b, &orderEvent)
	return orderEvent, err
}

func CreatePayment(createOrderEvent nozama.CreateOrderEvent) error {

	var paymentItem nozama.PaymentItem

	paymentItem.OrderID = createOrderEvent.OrderID
	paymentItem.TotalPrice = createOrderEvent.TotalPrice
	paymentItem.Status = nozama.PaymentStatusPending
	paymentItem.PaymentID = nozama.GeneratePrimaryKey()

	err := nozama.PutItem(paymentItem, nozama.PaymentsDynamoDBTableName)

	if err != nil {
		log.Printf("CreatePayment: %s", err)
		return err
	}

	return nil
}
