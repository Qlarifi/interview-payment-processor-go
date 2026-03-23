package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/socure/interview-payment-processor-go/model"
	"github.com/socure/interview-payment-processor-go/queue"
	"github.com/socure/interview-payment-processor-go/service"
)

const (
	defaultRecentLimit = 10
	maxRecentLimit     = 100
)

type PaymentHandler struct {
	queue   queue.PaymentMessageQueue
	service *service.PaymentService
}

func NewPaymentHandler(q queue.PaymentMessageQueue, svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{queue: q, service: svc}
}

func (h *PaymentHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/payments/submit", h.submit)
	mux.HandleFunc("GET /api/payments/recent", h.recent)
}

type submitRequest struct {
	PaymentID string  `json:"paymentId"`
	Amount    float64 `json:"amount"`
}

func (h *PaymentHandler) submit(w http.ResponseWriter, r *http.Request) {
	var req submitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	msg := model.NewPaymentMessage(req.PaymentID, req.Amount)
	h.queue.Submit(msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "accepted",
		"paymentId": req.PaymentID,
	})
}

func (h *PaymentHandler) recent(w http.ResponseWriter, r *http.Request) {
	limit := defaultRecentLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if limit < 1 {
		limit = 1
	}
	if limit > maxRecentLimit {
		limit = maxRecentLimit
	}

	payments := h.service.GetRecentPayments(limit)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}
