package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	nozama "nozama/src/shared"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handleUpdatePayment)
}

func handleUpdatePayment(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var httpErrorMessage = nozama.HttpMessageResponse{Message: "An error occured trying to update the payment"}

	if req.HTTPMethod != http.MethodPatch {
		return nozama.HttpResponse(http.StatusBadRequest, httpErrorMessage)
	}

	paymentRequest, err := toPaymentRequest(req.Body)

	if err != nil {
		log.Printf("NOZAMA - Cannot resolve payment request req %s with error %s", req.Body, err)
	}

	paymentEvent, err := updatePayment(paymentRequest)

	if err != nil {
		return nozama.HttpResponse(http.StatusBadRequest, httpErrorMessage)
	}

	err = sendPaymentEvent(paymentEvent)

	if err != nil {
		return nozama.HttpResponse(http.StatusBadRequest, httpErrorMessage)
	} else {
		return nozama.HttpResponse(http.StatusCreated, paymentEvent)
	}

}

func updatePayment(paymentRequest nozama.ProcessPaymentRequest) (nozama.UpdatedPaymentEvent, error) {

	paymentItem, err := nozama.GetPaymentByOrderID(paymentRequest.OrderID)

	if err != nil {
		log.Fatalf("An error ocurred while updating payment. Error: %s", err)
		return nozama.UpdatedPaymentEvent{}, err
	}

	log.Printf("NOZAMA - Moving Payment %s from %s to %s",
		paymentItem.PaymentID, paymentItem.Status, paymentRequest.Status)
	paymentItem.Status = paymentRequest.Status

	err = nozama.PutItem(paymentItem, nozama.PaymentsDynamoDBTableName)

	if err != nil {
		log.Fatalf("An error ocurred while updating payment. Error: %s", err)
		return nozama.UpdatedPaymentEvent{}, err
	}

	var newUpdatePaymentEvent nozama.UpdatedPaymentEvent
	newUpdatePaymentEvent.OrderID = paymentRequest.OrderID
	newUpdatePaymentEvent.Status = paymentRequest.Status

	return newUpdatePaymentEvent, nil

}

func sendPaymentEvent(newUpdatePaymentEvent nozama.UpdatedPaymentEvent) error {
	err := nozama.SendMessage(newUpdatePaymentEvent, nozama.OrdersSQSQueue)
	if err != nil {
		log.Printf("OnPaymentUpdatedException could not send UpdatePaymentEvent %s", err)
	}

	return err
}

func toPaymentRequest(body string) (nozama.ProcessPaymentRequest, error) {
	b := []byte(body)
	var paymentRequest nozama.ProcessPaymentRequest
	err := json.Unmarshal(b, &paymentRequest)

	if err != nil {
		log.Fatalf("An error ocurred while UnMarshal payment Body: %s. Error: %s",
			body, err)
		return nozama.ProcessPaymentRequest{}, err
	}

	if paymentRequest.OrderID == "" {
		err = errors.New("invalid payment orderID")
	}

	if paymentRequest.Status != nozama.PaymentStatusSuccess && paymentRequest.Status != nozama.PaymentStatusRejected {
		err = errors.New("invalid payment status")
	}

	return paymentRequest, err
}
