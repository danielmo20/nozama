package main

import (
	"encoding/json"
	nozama "nozama/src/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	nozama.Dynamo = &nozama.MockDynamoDB{}
}

func TestToPaymentRequest(t *testing.T) {
	paymentRequestExpected := nozama.ProcessPaymentRequest{
		OrderID: nozama.GeneratePrimaryKey(),
		Status:  nozama.PaymentStatusSuccess,
	}
	bytes, _ := json.Marshal(paymentRequestExpected)
	body := string(bytes)

	paymentRequest, err := ToPaymentRequest(body)

	assert.NoError(t, err)
	assert.EqualValues(t, paymentRequestExpected.OrderID, paymentRequest.OrderID)
	assert.EqualValues(t, paymentRequestExpected.Status, paymentRequest.Status)

}

func TestUpdateOrder_Success(t *testing.T) {
	paymentRequest := nozama.ProcessPaymentRequest{
		OrderID: nozama.GeneratePrimaryKey(),
		Status:  nozama.PaymentStatusSuccess,
	}

	paymentEvent, err := UpdatePayment(paymentRequest)

	assert.NoError(t, err)
	assert.EqualValues(t, paymentRequest.OrderID, paymentEvent.OrderID)
	assert.EqualValues(t, paymentRequest.Status, paymentEvent.Status)

}

func TestUpdateOrder_Rejected(t *testing.T) {
	paymentRequest := nozama.ProcessPaymentRequest{
		OrderID: nozama.GeneratePrimaryKey(),
		Status:  nozama.PaymentStatusRejected,
	}

	paymentEvent, err := UpdatePayment(paymentRequest)

	assert.NoError(t, err)
	assert.EqualValues(t, paymentRequest.OrderID, paymentEvent.OrderID)
	assert.EqualValues(t, paymentRequest.Status, paymentEvent.Status)

}
