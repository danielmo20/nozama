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

func handleUpdateOrder(ctx context.Context, sqsEvent events.SQSEvent) error {

	messageBody := sqsEvent.Records[0].Body

	log.Printf("handleUpdateOrder: SQS Event received: %s ", messageBody)

	updatePaymentEvent, err := ToUpdatePaymentEvent(messageBody)

	if err != nil {
		log.Printf("handleUpdateOrder: Error %s", err)
		return err
	}

	err = UpdateOrder(updatePaymentEvent)

	if err != nil {
		log.Printf("handleUpdateOrder: Error %s", err)
		return err
	}

	return err
}

func main() {
	lambda.Start(handleUpdateOrder)
}

func ToUpdatePaymentEvent(body string) (nozama.UpdatedPaymentEvent, error) {
	b := []byte(body)
	var paymentEvent nozama.UpdatedPaymentEvent
	err := json.Unmarshal(b, &paymentEvent)
	return paymentEvent, err
}

func UpdateOrder(updatePaymentEvent nozama.UpdatedPaymentEvent) error {

	orderItem, err := nozama.GetOrderByID(updatePaymentEvent.OrderID)

	if err != nil {
		log.Printf("UpdateOrder: An error ocurred while updating the order %s", err)
		return err
	}

	log.Printf("UpdateOrder: Moving Order %s from %s to paymentEvent.Status %s",
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
		log.Printf("UpdateOrder: An error ocurred while placing order. Error %s", err)
		return err
	}

	return nil
}
