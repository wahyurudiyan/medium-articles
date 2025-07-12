# Circuit Breaker Pattern in Go

This project provides a practical and easy-to-understand implementation of the Circuit Breaker pattern in Go. The Circuit Breaker is a crucial design pattern for building resilient distributed systems, as it helps prevent cascading failures and improves the stability of applications that depend on unreliable or slow external services.

## What is the Circuit Breaker Pattern?

The Circuit Breaker pattern acts as a protective barrier between your application and external services. When a service starts failing repeatedly, the circuit breaker "opens" to block further requests for a specified period, allowing the external service time to recover. This prevents your application from wasting resources on operations likely to fail and helps maintain overall system health.

### States

- **Closed:** Requests flow normally. Failures are monitored.
- **Open:** Requests are blocked immediately. The system waits for a cooldown interval.
- **Half-Open:** A limited number of test requests are allowed to check if the external service has recovered.

## Features

- **Configurable Parameters:** Set failure thresholds, cooldown intervals, and success rates to fit your use case.
- **Thread-Safe:** Safe for concurrent use in multi-goroutine environments.
- **Clear State Transitions:** Easily track and log transitions between Closed, Open, and Half-Open states.
- **Comprehensive Example:** Includes a sample `main.go` demonstrating real-world usage.
- **Unit Tests:** Thoroughly tested for reliability and correctness.

## Getting Started

### Prerequisites

- Go 1.18 or newer

### Installation

Clone the repository:

```sh
git clone https://github.com/yourusername/circuit-breaker-go.git
cd circuit-breaker-go
```

### Usage Example

Run the example program:

```sh
go run main.go
```

Sample output:

```
2025/07/08 21:46:16 INFO Trying to do work... state=Closed
2025/07/08 21:46:16 ERROR Attempt 1:  state=Closed error="simulated error"
2025/07/08 21:46:16 INFO Trying to do work... state=Closed
2025/07/08 21:46:16 ERROR Attempt 2:  state=Open error="simulated error"
2025/07/08 21:46:16 ERROR Attempt 3:  state=Open error="request blocked due to circuit open"
2025/07/08 21:46:16 INFO Waiting for interval to expire... state=Open
2025/07/08 21:46:19 INFO Attempt 4: success state=HalfOpen
2025/07/08 21:46:19 INFO Trying to do work after interval...
2025/07/08 21:46:19 INFO After interval state=Closed result=success
```

### How It Works

1. **Normal Operation (Closed):**  
    The circuit breaker allows all requests and monitors for failures.
2. **Failure Threshold Reached (Open):**  
    After a set number of consecutive failures, the breaker opens and blocks further requests for a cooldown interval.
3. **Recovery Check (Half-Open):**  
    After the interval, a limited request is allowed. If it succeeds, the breaker closes; if it fails, it reopens.

## Project Structure

- `breaker/`  
  Contains the core circuit breaker implementation and unit tests.
- `main.go`  
  Demonstrates how to use the circuit breaker in a real application.
- `README.md`  
  Project documentation.

## Running Tests

To run the unit tests:

```sh
cd breaker
go test
```

And latest coverate and race-condition test result:

```
  001_circuit_breaker_pattern on   main !3 ❯ go test -race -cover ./breaker
ok      github.com/wahyurudiyan/go-circuit-breaker/breaker      1.060s  coverage: 94.7% of statements
```

## When to Use

- Protecting your application from unreliable or slow external APIs.
- Preventing resource exhaustion due to repeated failures.
- Improving system stability and user experience in distributed systems.

## License

MIT

---

Feel free to contribute or open issues for suggestions and improvements!
