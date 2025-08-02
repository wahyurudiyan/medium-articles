package breaker

import (
	"context"
	"errors"
	"strings"
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
	if err == nil || !strings.Contains(err.Error(), "request blocked due to circuit open") {
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
		t.Fatalf("expected error, got %v", err)
	}
	if cb.state != StateHalfOpen {
		t.Fatalf("expected state half-open, got %v", cb.state)
	}
	if result == nil {
		t.Fatalf("expected not nil result")
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
func TestState_String(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{StateOpen, "Open"},
		{StateClosed, "Closed"},
		{StateHalfOpen, "HalfOpen"},
		{State(99), "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.state.String(); got != tt.expected {
			t.Errorf("State(%d).String() = %q, want %q", tt.state, got, tt.expected)
		}
	}
}

func TestCount_OnSuccess_OnFailure_Reset(t *testing.T) {
	c := &Count{}
	c.OnSuccess()
	if c.RequestSuccess != 1 || c.ConsecutiveSuccess != 1 || c.ConsecutiveFailure != 0 {
		t.Errorf("OnSuccess failed: %+v", c)
	}
	c.OnFailure()
	if c.RequestFailure != 1 || c.ConsecutiveFailure != 1 || c.ConsecutiveSuccess != 0 {
		t.Errorf("OnFailure failed: %+v", c)
	}
	c.Reset()
	if c.RequestSuccess != 0 || c.RequestFailure != 0 || c.ConsecutiveSuccess != 0 || c.ConsecutiveFailure != 0 {
		t.Errorf("Reset failed: %+v", c)
	}
}

func TestSetting_Setters(t *testing.T) {
	s := NewSetting().
		SetInterval(2 * time.Second).
		SetTimeout(5 * time.Second).
		SetFailureThreshold(10).
		SetSuccessRate(0.8)
	if s.Interval != 2*time.Second {
		t.Errorf("SetInterval failed")
	}
	if s.Timeout != 5*time.Second {
		t.Errorf("SetTimeout failed")
	}
	if s.FailureThreshold != 10 {
		t.Errorf("SetFailureThreshold failed")
	}
	if s.SuccessRateThreshold != 0.8 {
		t.Errorf("SetSuccessRate failed")
	}
	s.SetSuccessRate(2)
	if s.SuccessRateThreshold != 1 {
		t.Errorf("SetSuccessRate >1 failed")
	}
	s.SetSuccessRate(-1)
	if s.SuccessRateThreshold != 0.5 {
		t.Errorf("SetSuccessRate <0 failed")
	}
}

func TestCircuitBreaker_CurrentState(t *testing.T) {
	cb := NewCircuitBreaker(NewSetting())
	cb.state = StateOpen
	if cb.CurrentState() != "Open" {
		t.Errorf("CurrentState() = %q, want %q", cb.CurrentState(), "Open")
	}
	cb.state = StateClosed
	if cb.CurrentState() != "Closed" {
		t.Errorf("CurrentState() = %q, want %q", cb.CurrentState(), "Closed")
	}
	cb.state = StateHalfOpen
	if cb.CurrentState() != "HalfOpen" {
		t.Errorf("CurrentState() = %q, want %q", cb.CurrentState(), "HalfOpen")
	}
}

func TestCircuitBreaker_calculateSuccessRate(t *testing.T) {
	cb := NewCircuitBreaker(NewSetting())
	if cb.calculateSuccessRate() != 0 {
		t.Errorf("Expected 0 success rate with no requests")
	}
	cb.count.RequestSuccess = 3
	cb.count.RequestFailure = 1
	want := float32(3) / float32(4)
	got := cb.calculateSuccessRate()
	if got != want {
		t.Errorf("calculateSuccessRate() = %v, want %v", got, want)
	}
}

func TestCircuitBreaker_InvalidState(t *testing.T) {
	cb := NewCircuitBreaker(NewSetting())
	cb.state = State(99)
	handler := func(ctx context.Context) (any, error) { return nil, nil }
	_, err := cb.Fire(handler)
	if err == nil || err.Error() != "invalid circuit breaker state" {
		t.Errorf("Expected invalid circuit breaker state error, got %v", err)
	}
}
