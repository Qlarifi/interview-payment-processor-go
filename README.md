# Payment Processor

A Go application that processes payment messages from an in-memory queue using multiple worker goroutines. No external services are required: the queue and data store are in-memory.

## Prerequisites

- Go 1.26+

## Build

```bash
go build -o payment-processor .
```

## Run

```bash
go run .
```

The app will:

1. Start an in-memory data store for payments.
2. Start four worker goroutines that poll the queue and process messages.
3. Expose a REST API at `http://localhost:8080`.

## API

- **Submit a payment for processing** (async; workers will process it):

  ```bash
  curl -X POST http://localhost:8080/api/payments/submit \
    -H "Content-Type: application/json" \
    -d '{"paymentId": "99", "amount": 50.00}'
  ```

  Returns `202 Accepted` with `{"status":"accepted","paymentId":"99"}`.

- **Get the most recent payments** (ordered by creation time, newest first):

  ```bash
  curl -s http://localhost:8080/api/payments/recent | jq
  ```

  Optional query param: `limit` (default 10, max 100). Example: `curl -s "http://localhost:8080/api/payments/recent?limit=25" | jq`

  Returns `200 OK` with a JSON array of payment objects (`id`, `paymentId`, `amount`, `createdAt`).

## Configuration

Configuration is via environment variables:

| Variable | Description | Default |
|---|---|---|
| `QUEUE_VISIBILITY_TIMEOUT_SECONDS` | Visibility timeout for received messages | `30` |
| `WORKER_THREAD_COUNT` | Number of worker goroutines | `4` |
| `WORKER_POLL_INTERVAL_MS` | Delay between polls when queue is empty | `500` |
| `PORT` | HTTP listen port | `8080` |
