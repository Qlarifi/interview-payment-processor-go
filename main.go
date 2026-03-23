package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/socure/interview-payment-processor-go/consumer"
	"github.com/socure/interview-payment-processor-go/handler"
	"github.com/socure/interview-payment-processor-go/queue"
	"github.com/socure/interview-payment-processor-go/service"
	"github.com/socure/interview-payment-processor-go/store"
)

func main() {
	visibilityTimeout := envInt("QUEUE_VISIBILITY_TIMEOUT_SECONDS", 30)
	workerCount := envInt("WORKER_THREAD_COUNT", 4)
	pollIntervalMs := envInt("WORKER_POLL_INTERVAL_MS", 500)
	port := envString("PORT", "8080")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	paymentQueue := queue.NewInMemoryPaymentQueue(visibilityTimeout)
	paymentStore := store.NewPaymentStore()
	paymentService := service.NewPaymentService(paymentStore)

	queueConsumer := consumer.NewPaymentQueueConsumer(paymentQueue, paymentService, pollIntervalMs, workerCount)
	queueConsumer.Start(ctx)
	log.Printf("started %d worker goroutines (poll interval %dms, visibility timeout %ds)",
		workerCount, pollIntervalMs, visibilityTimeout)

	mux := http.NewServeMux()
	paymentHandler := handler.NewPaymentHandler(paymentQueue, paymentService)
	paymentHandler.Register(mux)

	srv := &http.Server{Addr: ":" + port, Handler: mux}
	log.Printf("listening on :%s", port)

	go func() {
		<-ctx.Done()
		log.Println("shutting down")
		srv.Shutdown(context.Background()) //nolint:errcheck
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
