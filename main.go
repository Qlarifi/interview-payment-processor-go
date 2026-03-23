package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/socure/interview-payment-processor-go/consumer"
	"github.com/socure/interview-payment-processor-go/handler"
	"github.com/socure/interview-payment-processor-go/queue"
	"github.com/socure/interview-payment-processor-go/service"
	"github.com/socure/interview-payment-processor-go/store"
)

func main() {
	visibilityTimeout := envInt("QUEUE_VISIBILITY_TIMEOUT_SECONDS", 30)
	threadCount := envInt("WORKER_THREAD_COUNT", 4)
	pollIntervalMs := envInt("WORKER_POLL_INTERVAL_MS", 500)
	port := envString("PORT", "8080")

	paymentQueue := queue.NewInMemoryPaymentQueue(visibilityTimeout)
	paymentStore := store.NewPaymentStore()
	paymentService := service.NewPaymentService(paymentStore)

	queueConsumer := consumer.NewPaymentQueueConsumer(paymentQueue, paymentService, pollIntervalMs, threadCount)
	queueConsumer.Start()
	log.Printf("started %d worker goroutines (poll interval %dms, visibility timeout %ds)",
		threadCount, pollIntervalMs, visibilityTimeout)

	mux := http.NewServeMux()
	paymentHandler := handler.NewPaymentHandler(paymentQueue, paymentService)
	paymentHandler.Register(mux)

	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func envInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envString(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
