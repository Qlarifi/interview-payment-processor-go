package worker

import (
	"context"
	"log"
	"time"

	"github.com/socure/interview-payment-processor-go/queue"
	"github.com/socure/interview-payment-processor-go/service"
)

type PaymentQueueWorker struct {
	queue          queue.PaymentMessageQueue
	service        *service.PaymentService
	pollIntervalMs int
	workerCount    int
}

func NewPaymentQueueWorker(
	q queue.PaymentMessageQueue,
	svc *service.PaymentService,
	pollIntervalMs int,
	workerCount int,
) *PaymentQueueWorker {
	return &PaymentQueueWorker{
		queue:          q,
		service:        svc,
		pollIntervalMs: pollIntervalMs,
		workerCount:    workerCount,
	}
}

// Start launches worker goroutines that poll the queue and process messages.
// Goroutines exit when ctx is cancelled.
func (w *PaymentQueueWorker) Start(ctx context.Context) {
	for i := 1; i <= w.workerCount; i++ {
		go w.pollLoop(ctx, i)
	}
}

func (w *PaymentQueueWorker) pollLoop(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d shutting down", workerID)
			return
		default:
		}

		msg := w.queue.Receive()
		if msg == nil {
			select {
			case <-ctx.Done():
				log.Printf("worker %d shutting down", workerID)
				return
			case <-time.After(time.Duration(w.pollIntervalMs) * time.Millisecond):
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
			w.service.ProcessPayment(*msg)
			w.queue.Ack(msg.MessageID)
		}()
	}
}
