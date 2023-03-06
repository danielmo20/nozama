package main

import (
	nozama "nozama/src/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToOrderRequest(t *testing.T) {

	body := "{\"user_id\": \"danielmo\",\"item\": \"test_item\",\"quantity\": 100,\"total_price\": 50000}"

	createOrderRequest, err := ToOrderRequest(body)

	assert.NoError(t, err)
	assert.EqualValues(t, "danielmo", createOrderRequest.UserID)
	assert.EqualValues(t, "test_item", createOrderRequest.Item)
	assert.EqualValues(t, 100, createOrderRequest.Quantity)
	assert.EqualValues(t, 50000, createOrderRequest.TotalPrice)

}

func TestCreateOrder(t *testing.T) {
	nozama.Dynamo = &nozama.MockDynamoDB{}

	orderRequest := nozama.CreateOrderRequest{
		UserID:     "danielmo",
		Item:       "test item",
		Quantity:   1,
		TotalPrice: 50000,
	}

	orderEvent, err := CreateOrder(orderRequest)

	assert.NotNil(t, orderEvent)
	assert.Nil(t, err)
	assert.NotEmpty(t, orderEvent.OrderID)
	assert.EqualValues(t, orderRequest.TotalPrice, orderEvent.TotalPrice)
}

func TestCreateOrder_DynamoError(t *testing.T) {
	nozama.Dynamo = &nozama.MockBrokenDynamoDB{}
	orderRequest := nozama.CreateOrderRequest{}
	orderEventResult, err := CreateOrder(orderRequest)
	orderEventExpected := nozama.CreateOrderEvent{}
	assert.Equal(t, orderEventExpected, orderEventResult)
	assert.Error(t, err)
}
