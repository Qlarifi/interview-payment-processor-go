package consumer

import (
	"log"
	"time"

	"github.com/socure/interview-payment-processor-go/queue"
	"github.com/socure/interview-payment-processor-go/service"
)

type PaymentQueueConsumer struct {
	queue          queue.PaymentMessageQueue
	service        *service.PaymentService
	pollIntervalMs int
	threadCount    int
}

func NewPaymentQueueConsumer(
	q queue.PaymentMessageQueue,
	svc *service.PaymentService,
	pollIntervalMs int,
	threadCount int,
) *PaymentQueueConsumer {
	return &PaymentQueueConsumer{
		queue:          q,
		service:        svc,
		pollIntervalMs: pollIntervalMs,
		threadCount:    threadCount,
	}
}

// Start launches worker goroutines that poll the queue and process messages.
func (c *PaymentQueueConsumer) Start() {
	for i := 1; i <= c.threadCount; i++ {
		go c.pollLoop(i)
	}
}

func (c *PaymentQueueConsumer) pollLoop(workerID int) {
	for {
		msg := c.queue.Receive()
		if msg == nil {
			time.Sleep(time.Duration(c.pollIntervalMs) * time.Millisecond)
			continue
		}
		log.Printf("worker %d processing payment %s", workerID, msg.PaymentID)
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("worker %d failed to process payment %s: %v", workerID, msg.PaymentID, r)
				}
			}()
			c.service.ProcessPayment(*msg)
			c.queue.Ack(msg.MessageID)
		}()
	}
}
