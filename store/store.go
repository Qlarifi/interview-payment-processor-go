package store

import (
	"sync"
	"time"

	"github.com/socure/interview-payment-processor-go/model"
)

// PaymentStore persists payment records in memory.
type PaymentStore struct {
	mu       sync.Mutex
	payments []model.Payment
	nextID   int64
}

func NewPaymentStore() *PaymentStore {
	return &PaymentStore{nextID: 1}
}

// Save persists a new payment and returns it with ID and CreatedAt populated.
func (s *PaymentStore) Save(paymentID string, amount float64) model.Payment {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := model.Payment{
		ID:        s.nextID,
		PaymentID: paymentID,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
	s.nextID++
	s.payments = append(s.payments, p)
	return p
}

// FindRecent returns up to limit payments ordered by creation time descending.
func (s *PaymentStore) FindRecent(limit int) []model.Payment {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := len(s.payments)
	if limit > n {
		limit = n
	}
	result := make([]model.Payment, limit)
	for i := range limit {
		result[i] = s.payments[n-1-i]
	}
	return result
}
