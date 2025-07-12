package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/wahyurudiyan/go-circuit-breaker/breaker"
)

func main() {
	cb := newCircuitBreaker(2, 2*time.Second)

	for i := 1; i <= 4; i++ {
		if i == 4 {
			logState("Waiting for interval to expire...", cb)
			time.Sleep(3 * time.Second)
		}
		fireAndLog(cb, i, alwaysFailHandler(cb))
	}

	// After interval, circuit breaker should allow a new attempt (half-open)
	fireAndLogResult(cb, "After interval", alwaysSucceedHandler())
}

func newCircuitBreaker(failureThreshold int, interval time.Duration) *breaker.CircuitBreaker {
	setting := breaker.NewSetting().
		SetFailureThreshold(failureThreshold).
		SetInterval(interval)
	return breaker.NewCircuitBreaker(setting)
}

func alwaysFailHandler(cb *breaker.CircuitBreaker) func(context.Context) (any, error) {
	return func(ctx context.Context) (any, error) {
		logState("Trying to do work...", cb)
		return nil, fmt.Errorf("simulated error")
	}
}

func alwaysSucceedHandler() func(context.Context) (any, error) {
	return func(ctx context.Context) (any, error) {
		slog.Info("Trying to do work after interval...")
		return "success", nil
	}
}

func fireAndLog(cb *breaker.CircuitBreaker, attempt int, handler func(context.Context) (any, error)) {
	_, err := cb.Fire(handler)
	if err != nil {
		slog.Error(fmt.Sprintf("Attempt %d: ", attempt), "state", cb.CurrentState(), "error", err)
	} else {
		slog.Info(fmt.Sprintf("Attempt %d: success", attempt), "state", cb.CurrentState())
	}
}

func fireAndLogResult(cb *breaker.CircuitBreaker, label string, handler func(context.Context) (any, error)) {
	result, err := cb.Fire(handler)
	if err != nil {
		slog.Error(label, "state", cb.CurrentState(), "error", err)
	} else {
		slog.Info(label, "state", cb.CurrentState(), "result", result)
	}
}

func logState(msg string, cb *breaker.CircuitBreaker) {
	slog.Info(msg, "state", cb.CurrentState())
}
