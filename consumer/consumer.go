package consumer

import (
	"context"
	"log"
	"time"

	"github.com/socure/interview-payment-processor-go/queue"
	"github.com/socure/interview-payment-processor-go/service"
)

type PaymentQueueConsumer struct {
	queue          queue.PaymentMessageQueue
	service        *service.PaymentService
	pollIntervalMs int
	workerCount    int
}

func NewPaymentQueueConsumer(
	q queue.PaymentMessageQueue,
	svc *service.PaymentService,
	pollIntervalMs int,
	workerCount int,
) *PaymentQueueConsumer {
	return &PaymentQueueConsumer{
		queue:          q,
		service:        svc,
		pollIntervalMs: pollIntervalMs,
		workerCount:    workerCount,
	}
}

// Start launches worker goroutines that poll the queue and process messages.
// Goroutines exit when ctx is cancelled.
func (c *PaymentQueueConsumer) Start(ctx context.Context) {
	for i := 1; i <= c.workerCount; i++ {
		go c.pollLoop(ctx, i)
	}
}

func (c *PaymentQueueConsumer) pollLoop(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d shutting down", workerID)
			return
		default:
		}

		msg := c.queue.Receive()
		if msg == nil {
			select {
			case <-ctx.Done():
				log.Printf("worker %d shutting down", workerID)
				return
			case <-time.After(time.Duration(c.pollIntervalMs) * time.Millisecond):
			}
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
