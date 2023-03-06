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

func TestToUpdatePaymentEvent(t *testing.T) {
	paymentEventExpected := nozama.UpdatedPaymentEvent{
		OrderID: nozama.GeneratePrimaryKey(),
		Status:  nozama.PaymentStatusSuccess,
	}
	bytes, _ := json.Marshal(paymentEventExpected)
	body := string(bytes)

	updatePaymentEvent, err := ToUpdatePaymentEvent(body)

	assert.NoError(t, err)
	assert.EqualValues(t, paymentEventExpected.OrderID, updatePaymentEvent.OrderID)
	assert.EqualValues(t, paymentEventExpected.Status, updatePaymentEvent.Status)

}

func TestUpdateOrder_Success(t *testing.T) {
	paymentEvent := nozama.UpdatedPaymentEvent{
		OrderID: nozama.GeneratePrimaryKey(),
		Status:  nozama.PaymentStatusSuccess,
	}

	err := UpdateOrder(paymentEvent)

	assert.NoError(t, err)

}

func TestUpdateOrder_Rejected(t *testing.T) {
	paymentEvent := nozama.UpdatedPaymentEvent{
		OrderID: nozama.GeneratePrimaryKey(),
		Status:  nozama.PaymentStatusRejected,
	}

	err := UpdateOrder(paymentEvent)

	assert.NoError(t, err)

}

func TestUpdateOrder_NotAllowed(t *testing.T) {
	paymentEvent := nozama.UpdatedPaymentEvent{
		OrderID: nozama.GeneratePrimaryKey(),
		Status:  nozama.PaymentStatusPending,
	}

	err := UpdateOrder(paymentEvent)

	assert.Error(t, err)

}
