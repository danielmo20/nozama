package nozama

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//const nozama_sqs_region = "us-east-2"
//const nozama_sqs_profile = "default"

var nozama_payments_sqs_queue = "SQSPayments"

//var nozama_payments_sqs_url = "https://sqs.us-east-2.amazonaws.com/055154340090/SQSPayments"

func SendMessage(messageBody string) error {

	log.Printf("NOZAMA - Sending message %s", messageBody)

	sqsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqsClient := sqs.New(sqsSession)

	sqsPaymentURL, err := GetQueueURL(nozama_payments_sqs_queue)

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
