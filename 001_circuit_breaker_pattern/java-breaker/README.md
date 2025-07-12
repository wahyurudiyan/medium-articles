# Circuit Breaker Pattern in Java

This project provides a practical and easy-to-understand implementation of the Circuit Breaker pattern in Java. The Circuit Breaker is a crucial design pattern for building resilient distributed systems, as it helps prevent cascading failures and improves the stability of applications that depend on unreliable or slow external services.

## What is the Circuit Breaker Pattern?

The Circuit Breaker pattern acts as a protective barrier between your application and external services. When a service starts failing repeatedly, the circuit breaker "opens" to block further requests for a specified period, allowing the external service time to recover. This prevents your application from wasting resources on operations likely to fail and helps maintain overall system health.

### States

- **Closed:** Requests flow normally. Failures are monitored.
- **Open:** Requests are blocked immediately. The system waits for a cooldown interval.
- **Half-Open:** A limited number of test requests are allowed to check if the external service has recovered.

## Features

- **Configurable Parameters:** Set failure thresholds, cooldown intervals, and success rates to fit your use case.
- **Thread-Safe:** Safe for concurrent use in multi-threaded environments.
- **Clear State Transitions:** Easily track and log transitions between Closed, Open, and Half-Open states.
- **Comprehensive Example:** Includes a sample `Main.java` demonstrating real-world usage.
- **Unit Tests:** Thoroughly tested for reliability and correctness.

## Getting Started

### Prerequisites

- Java 17 or newer
- Maven

### Installation

Clone the repository and enter the folder:

```sh
git clone https://github.com/yourusername/circuit-breaker-java.git
cd circuit-breaker-java
```

Build the project:

```sh
mvn clean package
```

### Usage Example

Run the example program:

```sh
java -cp target/java-breaker-1.0.0.jar com.resiliency.Main
```

Sample output:

```
=== START CIRCUIT BREAKER SIMULATION ===
[FAILURE] #1
State: CLOSE
[FAILURE] #2
State: CLOSE
[FAILURE] #3
State: OPEN

[INFO] Trying to access while state is OPEN...
acquirePermission(): false
State: OPEN

[INFO] Waiting 2000ms to enter HALF_OPEN...
[INFO] Trying to access again after interval...
acquirePermission(): true
State: HALF_OPEN
[SUCCESS] #1
State: HALF_OPEN
[SUCCESS] #2
State: CLOSE

[INFO] Final status:
State: CLOSE
=== END ===
```

## How It Works

1. **Normal Operation (Closed):**  
    The circuit breaker allows all requests and monitors for failures.
2. **Failure Threshold Reached (Open):**  
    After a set number of consecutive failures, the breaker opens and blocks further requests for a cooldown interval.
3. **Recovery Check (Half-Open):**  
    After the interval, a limited request is allowed. If it succeeds, the breaker closes; if it fails, it reopens.

## Project Structure

- `src/main/java/breaker/`  
  Contains the core circuit breaker implementation.
- `src/main/java/com/resiliency/Main.java`  
  Demonstrates how to use the circuit breaker in a real application.
- `src/test/java/breaker/`  
  Unit tests for the circuit breaker.
- `README.md`  
  Project documentation.

## Running Tests

To run the unit tests:

```sh
mvn test
```

## When to Use

- Protecting your application from unreliable or slow external APIs.
- Preventing resource exhaustion due to repeated failures.
- Improving system stability and user experience in distributed systems.

## License

MIT

---

Feel free to contribute or open issues for suggestions and improvements!
