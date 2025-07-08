# Circuit Breaker Pattern in Go

This project demonstrates a simple implementation of the Circuit Breaker pattern in Go. The circuit breaker helps prevent repeated failures when calling unreliable services, by blocking calls for a period after a threshold of failures is reached.

## Features
- Configurable failure threshold, interval, and success rate
- Thread-safe state transitions (Closed, Open, Half-Open)
- Example usage in `main.go`
- Unit tests included

## Usage Example

```
go run main.go
```

Example output:
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

## How It Works
- The circuit starts in the **Closed** state, allowing requests.
- After a configurable number of consecutive failures, it transitions to **Open** and blocks further requests for a set interval.
- After the interval, it transitions to **Half-Open** and allows a test request.
- If the test request succeeds, the circuit closes; if it fails, it reopens.

## Project Structure
- `breaker/` - Contains the circuit breaker implementation and tests
- `main.go` - Example usage
- `README.md` - This file

## Running Tests

```
cd breaker

go test
```

## License
MIT
