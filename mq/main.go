package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/remind101/mq-go"

	uuid "github.com/gofrs/uuid"
)

const messages = 1000

var ops uint64
var wg sync.WaitGroup

func getID() *string {
	u4, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("could not generate uuid [%+v]\n", err)
	}

	r := u4.String()
	return &r
}

func main() {
	queueURL := "https://sqs.us-east-1.amazonaws.com/.../mq.fifo"

	h := mq.HandlerFunc(func(m *mq.Message) error {
		defer wg.Done()

		fmt.Printf("Received message: %s\n", aws.StringValue(m.SQSMessage.Body))

		atomic.AddUint64(&ops, 1)

		// Returning no error signifies the message was processed successfully.
		// The Server will queue the message for deletion.
		return nil
	})

	ctx := context.Background()
	// Configure mq.Server
	s := mq.NewServer(queueURL, h)

	// Start a loop to receive SQS messages and pass them to the Handler.
	s.Start()
	defer s.Shutdown(ctx)

	// Handle SIGINT and SIGTERM gracefully.
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh

		fmt.Printf("SIGINT|SIGTERM received: %v\n", sig)

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// We received an interrupt signal, shut down gracefully.
		if err := s.Shutdown(ctx); err != nil {
			fmt.Printf("SQS server shutdown: %v\n", err)
		}

		os.Exit(1)
	}()

	// Start a publisher
	p := mq.NewPublisher(queueURL)
	p.Start()
	defer p.Shutdown(ctx)

	wg.Add(1)

	go func() {
		defer wg.Done()
		max := 1000 + messages
		for i := 1001; i <= max; i++ {
			idx := i - 1000
			wg.Add(1)
			// Publish messages (will be batched).
			p.Publish(&sqs.SendMessageBatchRequestEntry{
				Id:             getID(),
				MessageGroupId: aws.String("Hello"),
				MessageBody:    aws.String(fmt.Sprintf("%v=>Hello-Message-%s", idx, *getID())),
			})
			log.Printf("published message %v\n", idx)
		}
	}()

	wg.Wait()

	fmt.Println("ops:", ops)
}
