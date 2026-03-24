package queue

import (
	"sync"

	"github.com/socure/interview-payment-processor-go/model"
)

// PaymentMessageQueue is the abstraction for the payment message queue.
// Implementations support submit, receive with visibility semantics, and ack.
type PaymentMessageQueue interface {
	// Submit adds a message to the queue for processing.
	Submit(message model.PaymentMessage)

	// Receive returns the next available message. The message should not be
	// visible to other workers until the visibility timeout expires or it
	// is acknowledged.
	// Returns nil if no message is available.
	Receive() *model.PaymentMessage

	// Ack acknowledges successful processing. The message should be removed
	// so it is not delivered again.
	Ack(messageID string)
}

// InMemoryPaymentQueue is an in-memory implementation of the payment message queue.
// Messages are held in a single pending queue. Receive returns the next message
// and the implementation is expected to hide it from other workers until
// the visibility timeout expires or the message is acknowledged.
type InMemoryPaymentQueue struct {
	mu                       sync.Mutex
	pending                  []model.PaymentMessage
	visibilityTimeoutSeconds int
}

// NewInMemoryPaymentQueue creates a new InMemoryPaymentQueue with the given
// visibility timeout in seconds.
func NewInMemoryPaymentQueue(visibilityTimeoutSeconds int) *InMemoryPaymentQueue {
	return &InMemoryPaymentQueue{
		visibilityTimeoutSeconds: visibilityTimeoutSeconds,
	}
}

func (q *InMemoryPaymentQueue) Submit(message model.PaymentMessage) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.pending = append(q.pending, message)
}

func (q *InMemoryPaymentQueue) Receive() *model.PaymentMessage {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.pending) == 0 {
		return nil
	}
	msg := q.pending[0]
	return &msg
}

func (q *InMemoryPaymentQueue) Ack(messageID string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i, m := range q.pending {
		if m.MessageID == messageID {
			q.pending = append(q.pending[:i], q.pending[i+1:]...)
			return
		}
	}
}
