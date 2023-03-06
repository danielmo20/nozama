package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	nozama "nozama/src/shared"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	messageBody := sqsEvent.Records[0].Body

	log.Printf("NOZAMA - event recieved: %s ", messageBody)

	updatePaymentEvent, err := toUpdatePaymentEvent(messageBody)

	if err != nil {
		log.Fatalf("An error while parsing toUpdatePaymentEvent ")
		return err
	}

	err = updateOrder(updatePaymentEvent)

	if err != nil {
		log.Fatalf("An error while creating a Payment")
		return err
	}

	return err
}

func main() {
	lambda.Start(handler)
}

func toUpdatePaymentEvent(body string) (nozama.UpdatedPaymentEvent, error) {
	b := []byte(body)
	var paymentEvent nozama.UpdatedPaymentEvent
	err := json.Unmarshal(b, &paymentEvent)
	return paymentEvent, err
}

func updateOrder(updatePaymentEvent nozama.UpdatedPaymentEvent) error {

	orderItem, err := nozama.GetOrderByID(updatePaymentEvent.OrderID)

	if err != nil {
		log.Fatalf("NOZAMA - An error ocurred while updating the order %s", err)
		return err
	}

	log.Printf("NOZAMA - Moving Order %s from %d to paymentEvent.Status %d",
		orderItem.OrderID, orderItem.Status, updatePaymentEvent.Status)
	switch updatePaymentEvent.Status {
	case nozama.PaymentStatusSuccess:
		{
			orderItem.Status = nozama.OrderStatusRedyForShipping
			break
		}
	case nozama.PaymentStatusRejected:
		{
			orderItem.Status = nozama.OrderStatusRejected
			break
		}
	default:
		{
			return errors.New("payment status logic not implemented yet")
		}
	}

	err = nozama.PutItem(orderItem, nozama.OrdersDynamoDBTableName)

	if err != nil {
		log.Fatalf("NOZAMA - An error ocurred while placing order %s", err)
		return err
	}

	return nil
}
