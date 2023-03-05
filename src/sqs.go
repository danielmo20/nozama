package nozama

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func SendMessage(messageBody string, queueName string) error {

	log.Printf("NOZAMA - Sending message %s", messageBody)

	sqsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqsClient := sqs.New(sqsSession)

	sqsPaymentURL, err := GetQueueURL(queueName)

	if err != nil {
		log.Printf("Failed to initialize new session: %v", err)
		return err
	}

	result, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &sqsPaymentURL,
		MessageBody: aws.String(messageBody),
	})

	if result == nil {
		log.Fatal("NOZAMA - Message not sent!")
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
		log.Fatalln(err)
	}

	return *qurlOut.QueueUrl, err
}
