package main

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	maxMessages                  = 1
	visibilityTimeout            = 30
	waitTimeSeconds              = 20
	visibilityTimeoutAfterFailed = 0
	cancelTimeout                = 11
)

var qURLv2 = "https://sqs.eu-west-1.amazonaws.com/041936244769/csqsv2.fifo"
var qURLv2Errors = "https://sqs.eu-west-1.amazonaws.com/041936244769/csqsv2errors.fifo"

type payloadMessage struct {
	message     *sqs.Message
	payloadType int // 0: message, 1: cancel
}

func processMessage(msg *sqs.Message) error {
	if strings.HasPrefix(*msg.Body, "ERROR:") {
		return errors.New("fake error message")
	}
	log.Printf("\n-----BEGIN MESSAGE-----\n%+v\n-----END MESSAGE-----\n", *msg.Body)
	return nil
}

func ackMessage(ctx context.Context, svc *sqs.SQS, msg *sqs.Message) error {
	_, err := svc.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &qURLv2,
		ReceiptHandle: msg.ReceiptHandle,
	})
	return err
}

func sendMessageError(ctx context.Context, svc *sqs.SQS, msg *sqs.Message) (*sqs.SendMessageOutput, error) {
	return svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageAttributes: msg.MessageAttributes,
		MessageBody:       msg.Body,
		MessageGroupId:    msg.Attributes["MessageGroupId"],
		QueueUrl:          &qURLv2Errors,
	})
}

func check(wg *sync.WaitGroup) {
	defer wg.Done()

	tt := []int64{1, 2, 3}

	for _, t := range tt {
		log.Printf("checking db status: %v\n", t)
		time.Sleep(time.Duration(t) * time.Second) // request health-check db
	}
}

func receiveMessages(ctx context.Context, wg *sync.WaitGroup, svc *sqs.SQS, chn chan<- *payloadMessage) {
	defer wg.Done()

	for {
		log.Printf("waiting for messages...\n")
		result, err := svc.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
				aws.String(sqs.MessageSystemAttributeNameMessageGroupId),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            &qURLv2,
			MaxNumberOfMessages: aws.Int64(maxMessages),
			VisibilityTimeout:   aws.Int64(visibilityTimeout),
			WaitTimeSeconds:     aws.Int64(waitTimeSeconds),
		})
		if err != nil {
			log.Println("receive messages error:", err)
			return
		}

		for _, message := range result.Messages {
			chn <- &payloadMessage{
				message:     message,
				payloadType: 0,
			}
		}

		for _, msg := range result.Messages {
			if err := processMessage(msg); err == nil {
				err := ackMessage(ctx, svc, msg)
				if err != nil {
					log.Printf("Ack error: %+v\n", err)
					return
				}
				log.Printf("Success: %+v\n", msg)
			} else {
				var failedVisibilityTimeout int64 = 0
				_, e := svc.ChangeMessageVisibilityWithContext(ctx, &sqs.ChangeMessageVisibilityInput{
					QueueUrl:          &qURLv2,
					ReceiptHandle:     msg.ReceiptHandle,
					VisibilityTimeout: &failedVisibilityTimeout,
				})
				if e != nil {
					log.Printf("returned to queue in %v seconds\n", visibilityTimeout)
				} else {
					log.Printf("returned to queue %v\n", qURLv2)
				}

				wg := sync.WaitGroup{}
				wg.Add(1)
				go check(&wg)
				wg.Wait()

				log.Printf("consumer restarted!\n")
			}
		}
	}
}

func consumeMessages(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, svc *sqs.SQS, chn chan *payloadMessage) {
	defer wg.Done()

	for {
		select {
		case payload := <-chn:
			if payload.payloadType == 1 {
				cancel()
				return
			}
			msg := payload.message
			if err := processMessage(msg); err == nil {
				err := ackMessage(ctx, svc, msg)
				if err != nil {
					log.Printf("Ack error: %+v\n", err)
					return
				}
				log.Printf("Success: %+v\n", msg)
			} else {
				var failedVisibilityTimeout int64 = 0
				_, e := svc.ChangeMessageVisibilityWithContext(ctx, &sqs.ChangeMessageVisibilityInput{
					QueueUrl:          &qURLv2,
					ReceiptHandle:     msg.ReceiptHandle,
					VisibilityTimeout: &failedVisibilityTimeout,
				})
				if e != nil {
					log.Printf("returned to queue in %v seconds\n", visibilityTimeout)
				} else {
					log.Printf("returned to queue %v\n", qURLv2)
				}

				wg := sync.WaitGroup{}
				wg.Add(1)
				go check(&wg)
				wg.Wait()

				log.Printf("consumer restarted!\n")
			}
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	chn := make(chan *payloadMessage)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go receiveMessages(ctx, &wg, svc, chn)

	wg.Add(1)
	go consumeMessages(ctx, cancel, &wg, svc, chn)

	wg.Wait()
}
