package model

import (
	"crypto/rand"
	"fmt"
	"time"
)

// PaymentMessage represents a single unit of work consumed from the queue by a worker.
type PaymentMessage struct {
	MessageID string  `json:"messageId"`
	PaymentID string  `json:"paymentId"`
	Amount    float64 `json:"amount"`
}

// NewPaymentMessage creates a PaymentMessage with a random message ID.
func NewPaymentMessage(paymentID string, amount float64) PaymentMessage {
	return PaymentMessage{
		MessageID: newUUID(),
		PaymentID: paymentID,
		Amount:    amount,
	}
}

// Payment represents a persisted payment record.
type Payment struct {
	ID        int64     `json:"id"`
	PaymentID string    `json:"paymentId"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
