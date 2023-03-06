package main

import (
	nozama "nozama/src/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCreateOrderEvent(t *testing.T) {

	body := "{\"order_id\": \"591c48b4-bee6-4885-b8f7-17500fa2cd31\",\"total_price\": 786400}"

	createOrderEvent, err := ToCreateOrderEvent(body)

	assert.NoError(t, err)
	assert.EqualValues(t, "591c48b4-bee6-4885-b8f7-17500fa2cd31", createOrderEvent.OrderID)
	assert.EqualValues(t, 786400, createOrderEvent.TotalPrice)

}

func TestCreatePayment(t *testing.T) {
	nozama.Dynamo = &nozama.MockDynamoDB{}

	createOrderEvent := nozama.CreateOrderEvent{
		OrderID:    nozama.GeneratePrimaryKey(),
		TotalPrice: 50000,
	}

	err := CreatePayment(createOrderEvent)

	assert.Nil(t, err)
}
