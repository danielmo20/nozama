package main

import (
	nozama "nozama/src/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpResponse(t *testing.T) {

	var newUpdatePaymentEvent nozama.UpdatedPaymentEvent
	newUpdatePaymentEvent.OrderID = "dummy-order-id"
	newUpdatePaymentEvent.Status = "dummy-payment-status"

	response, err := nozama.HttpResponse(201, newUpdatePaymentEvent)

	assert.NoError(t, err)
	assert.EqualValues(t, 201, response.StatusCode)
	assert.Contains(t, response.Body, newUpdatePaymentEvent.OrderID)
	assert.Contains(t, response.Body, "order_id")
	assert.EqualValues(t, "application/json", response.Headers["Content-Type"])

}
