package service

import (
	"log"

	"github.com/socure/interview-payment-processor-go/model"
	"github.com/socure/interview-payment-processor-go/store"
)

type PaymentService struct {
	store *store.PaymentStore
}

func NewPaymentService(store *store.PaymentStore) *PaymentService {
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
