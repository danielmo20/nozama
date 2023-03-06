package nozama

import (
	nozama "nozama/src/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	nozama.Dynamo = &nozama.MockDynamoDB{}
}

func TestGeneratePrimaryKey(t *testing.T) {
	primaryKey := nozama.GeneratePrimaryKey()

	assert.NotEmpty(t, primaryKey)
}

func TestGetOrderByID(t *testing.T) {
	orderId := nozama.GeneratePrimaryKey()
	order, err := nozama.GetOrderByID(orderId)

	assert.NoError(t, err)
	assert.EqualValues(t, orderId, order.OrderID)
}

func TestGetPaymentByOrderID(t *testing.T) {
	orderId := nozama.GeneratePrimaryKey()
	paymentItem, err := nozama.GetPaymentByOrderID(orderId)

	assert.NoError(t, err)
	assert.EqualValues(t, orderId, paymentItem.OrderID)
}

func TestPutItemOrder(t *testing.T) {
	orderItem := nozama.MockOrderIncomplete(nozama.GeneratePrimaryKey())
	err := nozama.PutItem(orderItem, nozama.OrdersDynamoDBTableName)

	assert.Nil(t, err)
}

func TestPutItemPayment(t *testing.T) {
	paymentItem := nozama.MockPaymentPending(nozama.GeneratePrimaryKey())
	err := nozama.PutItem(paymentItem, nozama.PaymentsDynamoDBTableName)

	assert.Nil(t, err)
}
