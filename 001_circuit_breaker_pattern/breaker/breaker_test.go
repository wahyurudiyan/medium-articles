package breaker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_ClosedState_Success(t *testing.T) {
	cb := NewCircuitBreaker(NewSetting())
	cb.state = StateClosed

	handler := func(ctx context.Context) (any, error) {
		return "ok", nil
	}

	result, err := cb.Fire(handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != "ok" {
		t.Fatalf("expected result 'ok', got %v", result)
	}
}

func TestCircuitBreaker_ClosedState_FailureToOpen(t *testing.T) {
	cb := NewCircuitBreaker(NewSetting().SetFailureThreshold(2))
	cb.state = StateClosed

	handler := func(ctx context.Context) (any, error) {
		return nil, errors.New("fail")
	}

	// First failure
	_, err := cb.Fire(handler)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if cb.state != StateClosed {
		t.Fatalf("expected state closed, got %v", cb.state)
	}

	// Second failure triggers open
	_, err = cb.Fire(handler)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if cb.state != StateOpen {
		t.Fatalf("expected state open, got %v", cb.state)
	}
}

func TestCircuitBreaker_OpenState_Blocks(t *testing.T) {
	cb := NewCircuitBreaker(NewSetting())
	cb.state = StateOpen
	cb.lastFailure = time.Now()

	handler := func(ctx context.Context) (any, error) {
		return "should not run", nil
	}

	result, err := cb.Fire(handler)
	if err == nil || err.Error() != "request blocked due to circuit open" {
		t.Fatalf("expected open state error, got %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestCircuitBreaker_OpenState_ToHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(NewSetting().SetInterval(10 * time.Millisecond))
	cb.state = StateOpen
	cb.lastFailure = time.Now().Add(-20 * time.Millisecond)

	handler := func(ctx context.Context) (any, error) {
		return "should not run", nil
	}

	result, err := cb.Fire(handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cb.state != StateHalfOpen {
		t.Fatalf("expected state half-open, got %v", cb.state)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestCircuitBreaker_HalfOpen_SuccessToClosed(t *testing.T) {
	cb := NewCircuitBreaker(NewSetting().SetSuccessRate(0.5))
	cb.state = StateHalfOpen
	cb.count.RequestSuccess = 3
	cb.count.RequestFailure = 1
	cb.lastFailure = time.Now().Add(-2 * time.Second)

	handler := func(ctx context.Context) (any, error) {
		return "ok", nil
	}

	result, err := cb.Fire(handler)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cb.state != StateClosed {
		t.Fatalf("expected state closed, got %v", cb.state)
	}
	if result != "ok" {
		t.Fatalf("expected result 'ok', got %v", result)
	}
}

func TestCircuitBreaker_HalfOpen_FailureStaysHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(NewSetting().SetSuccessRate(0.9))
	cb.state = StateHalfOpen
	cb.count.RequestSuccess = 1
	cb.count.RequestFailure = 1
	cb.lastFailure = time.Now().Add(-2 * time.Second)

	handler := func(ctx context.Context) (any, error) {
		return nil, errors.New("fail")
	}

	result, err := cb.Fire(handler)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if cb.state != StateHalfOpen {
		t.Fatalf("expected state half-open, got %v", cb.state)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}
