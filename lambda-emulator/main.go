package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"fmt"
)

var (
	sqsSvc   *sqs.SQS
	endpoint string = "http://sqs:9324"
	queueURL string = "http://sqs:9324/queue/messages"
)

func pollMessages(chn chan<- *sqs.Message) {
	fmt.Println("asked to poll")
	for {
		output, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: aws.Int64(2),
			WaitTimeSeconds:     aws.Int64(5),
		})
		fmt.Println("finished polling")
		if err != nil {
			fmt.Println("failed to fetch sqs message", err)
		}

		for _, message := range output.Messages {
			fmt.Println(message)
			chn <- message
		}

	}

}

// pullMessage...
func pullMessage(msg *sqs.Message) string {

	fmt.Println("message received", *msg.MessageId)

	return *msg.Body
}

func deleteMessage(msg *sqs.Message) {
	sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	fmt.Println("message deleted", *msg.MessageId)
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println(err)
	}
	sqsSvc = sqs.New(sess, aws.NewConfig().WithEndpoint(endpoint))
	fmt.Println("polling")
	messages := make(chan *sqs.Message, 2)
	go pollMessages(messages)
	for message := range messages {
		msg := pullMessage(message)
		fmt.Println("message pulled", msg)
		deleteMessage(message)
	}
	fmt.Println("not polling")
}
