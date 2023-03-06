package main

import (
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
