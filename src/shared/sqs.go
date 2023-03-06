package nozama

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func SendMessage(event interface{}, queueName string) error {

	messageBytes, err := json.Marshal(event)

	if err != nil {
		log.Printf("SendMessage: Cannont json.Marshal this %v", event)
	}

	messageBody := string(messageBytes)

	log.Printf("SendMessage: Sending message %s", messageBody)

	sqsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqsClient := sqs.New(sqsSession)

	sqsPaymentURL, err := GetQueueURL(queueName)

	if err != nil {
		log.Printf("SendMessage: failed to initialize new session: %v", err)
		return err
	}

	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &sqsPaymentURL,
		MessageBody: aws.String(messageBody),
	})

	if err != nil {
		log.Printf("SendMessage: Message not sent! Error: %s", err)
	}

	return err
}

func GetQueueURL(queue string) (string, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	qurlOut, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		log.Printf("GetQueueURL: %s", err)
		return "", err
	}

	return *qurlOut.QueueUrl, err
}
