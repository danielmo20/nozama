package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var nozama_orders_sqs_queue = "SQSOrders"

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("NOZAMA - The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
