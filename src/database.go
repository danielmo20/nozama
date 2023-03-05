package nozama

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

const PaymentsTableName = "payments"

var dynamo *dynamodb.DynamoDB

func init() {
	dynamo = connect()
}

func connect() (db *dynamodb.DynamoDB) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
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

	_, err = dynamo.PutItem(input)
	return err
}

func GeneratePrimaryKey() string {
	return uuid.NewString()
}
