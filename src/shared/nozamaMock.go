package nozama

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (dbmock *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {

	itemDetail := MockOrderIncomplete(*input.Key["order_id"].S)

	result := &dynamodb.GetItemOutput{
		Item: itemDetail,
	}

	return result, nil

}

func (dbmock *MockDynamoDB) Query(query *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {

	itemDetail := MockPaymentPending(*query.ExpressionAttributeValues[":orderId"].S)

	result := &dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			itemDetail,
		},
	}

	return result, nil

}

func (dbmock *MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, nil
}

type MockBrokenDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (dbmock *MockBrokenDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, errors.New("this a mock error while saving the record")
}

func MockOrderIncomplete(orderID string) map[string]*dynamodb.AttributeValue {
	return MockOrderGetItemOutput(orderID, OrderStatusIncomplete)
}

func MockOrderReadyForShipping(orderID string) map[string]*dynamodb.AttributeValue {
	return MockOrderGetItemOutput(orderID, OrderStatusRedyForShipping)
}

func MockOrderRejected(orderID string) map[string]*dynamodb.AttributeValue {
	return MockOrderGetItemOutput(orderID, OrderStatusRejected)
}

func MockPaymentPending(orderID string) map[string]*dynamodb.AttributeValue {
	return MockPaymentGetItemOutput(orderID, PaymentStatusPending)
}

func MockOrderGetItemOutput(orderId, status string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"order_id": {
			S: aws.String(orderId),
		},
		"user_id": {
			S: aws.String("danielmo"),
		},
		"item": {
			S: aws.String("test item"),
		},
		"quantity": {
			N: aws.String("1"),
		},
		"total_price": {
			N: aws.String("50000"),
		},
		"status": {
			N: aws.String(status),
		},
	}
}

func MockPaymentGetItemOutput(orderId, status string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"payment_id": {
			S: aws.String(GeneratePrimaryKey()),
		},
		"order_id": {
			S: aws.String(orderId),
		},
		"total_price": {
			N: aws.String("50000"),
		},
		"status": {
			S: aws.String(status),
		},
	}
}
