package service

import (
	"log"

	"github.com/socure/interview-payment-processor-go/model"
)

type PaymentStore interface {
	Save(paymentID string, amount float64) model.Payment
	FindRecent(limit int) []model.Payment
}

type PaymentService struct {
	store PaymentStore
}

func NewPaymentService(store PaymentStore) *PaymentService {
	return &PaymentService{store: store}
}

func (s *PaymentService) ProcessPayment(msg model.PaymentMessage) model.Payment {
	p := s.store.Save(msg.PaymentID, msg.Amount)
	log.Printf("payment %s created -- id %d", msg.PaymentID, p.ID)
	return p
}

func (s *PaymentService) GetRecentPayments(limit int) []model.Payment {
	return s.store.FindRecent(limit)
}
