package main

import (
	"context"
	"encoding/json"
	"log"
	nozama "nozama/src/shared"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	messageBody := sqsEvent.Records[0].Body

	log.Printf("NOZAMA - event recieved: %s ", messageBody)

	createOrderEvent, err := toCreateOrderEvent(messageBody)

	if err != nil {
		log.Fatalf("An error while parsing toCreateOrderEvent ")
		return err
	}

	err = createPayment(createOrderEvent)

	if err != nil {
		log.Fatalf("An error while creating a Payment")
		return err
	}

	return err
}

func main() {
	lambda.Start(handler)
}

func toCreateOrderEvent(body string) (nozama.CreateOrderEvent, error) {
	b := []byte(body)
	var orderEvent nozama.CreateOrderEvent
	err := json.Unmarshal(b, &orderEvent)
	return orderEvent, err
}

func createPayment(createOrderEvent nozama.CreateOrderEvent) error {

	var paymentItem nozama.PaymentItem

	paymentItem.OrderID = createOrderEvent.OrderID
	paymentItem.TotalPrice = createOrderEvent.TotalPrice
	paymentItem.Status = nozama.PaymentStatusPending
	paymentItem.PaymentID = nozama.GeneratePrimaryKey()

	err := nozama.PutItem(paymentItem, nozama.PaymentsDynamoDBTableName)

	if err != nil {
		log.Fatalf("An error ocurred while placing order %s", err)
		return err
	}

	return nil
}
