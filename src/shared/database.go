package nozama

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

var Dynamo dynamodbiface.DynamoDBAPI

func init() {
	Dynamo = connect()
}

func connect() (db *dynamodb.DynamoDB) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}

func GetPaymentByOrderID(orderID string) (PaymentItem, error) {

	_, err := uuid.Parse(orderID)

	if err != nil {
		log.Printf("GetPaymentByOrderID: orderID %s invalid %s", orderID, err)
		return PaymentItem{}, err
	}

	result, err := Dynamo.Query(&dynamodb.QueryInput{
		TableName: aws.String(PaymentsDynamoDBTableName),
		IndexName: aws.String("order_id-index"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":orderId": {
				S: aws.String(orderID),
			},
		},
		KeyConditionExpression: aws.String("order_id = :orderId"),
	})

	if err != nil {
		log.Printf("GetPaymentByOrderID: orderId %s. Error: %s", orderID, err)
		return PaymentItem{}, err
	}

	if result.Items == nil && result.Items[0] == nil {
		log.Printf("GetPaymentByOrderID: orderID %s doesn't exist", orderID)
		return PaymentItem{}, err
	}

	paymentItem := PaymentItem{}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &paymentItem)

	if err != nil {
		log.Printf("GetPaymentByOrderID: an error ocurred while parsing PaymentItem %v. Error: %s", paymentItem, err)
		return PaymentItem{}, err
	}

	return paymentItem, nil
}

func GetOrderByID(orderID string) (OrderItem, error) {

	_, err := uuid.Parse(orderID)

	if err != nil {
		log.Printf("GetOrderByID: orderID %s invalid. %s", orderID, err)
		return OrderItem{}, err
	}

	result, err := Dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(OrdersDynamoDBTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"order_id": {
				S: aws.String(orderID),
			},
		},
	})

	if err != nil {
		log.Printf("GetOrderByID: %s. Error: %s", orderID, err)
		return OrderItem{}, err
	}

	if result.Item == nil {
		log.Printf("GetOrderByID: orderID %s doesn't exist", orderID)
		return OrderItem{}, err
	}

	orderItem := OrderItem{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &orderItem)

	if err != nil {
		log.Printf("GetOrderByID: An error ocurred while parsing OrderItem %v. Error: %s", orderItem, err)
		return OrderItem{}, err
	}

	return orderItem, nil
}

func PutItem(newRecord interface{}, tableName string) error {
	entity, err := dynamodbattribute.MarshalMap(newRecord)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      entity,
		TableName: aws.String(tableName),
	}

	_, err = Dynamo.PutItem(input)
	return err
}

// Generates a UUID using the Google-uuid import
func GeneratePrimaryKey() string {
	return uuid.NewString()
}
