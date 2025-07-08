package breaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

type State int8

const (
	StateOpen State = iota + 1
	StateClosed
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateOpen:
		return "Open"
	case StateClosed:
		return "Closed"
	case StateHalfOpen:
		return "HalfOpen"
	default:
		return "Unknown"
	}
}

const (
	defaultTimeout              = 30 * time.Second
	defaultInterval             = 1 * time.Second
	defaultFailureThreshold     = 5
	defaultSuccessRateThreshold = 0.5
)

type HandlerFunc func(ctx context.Context) (any, error)

type Count struct {
	RequestSuccess     int
	RequestFailure     int
	ConsecutiveSuccess int
	ConsecutiveFailure int
}

func (c *Count) OnSuccess() {
	c.RequestSuccess++
	c.ConsecutiveSuccess++
	c.ConsecutiveFailure = 0
}

func (c *Count) OnFailure() {
	c.RequestFailure++
	c.ConsecutiveFailure++
	c.ConsecutiveSuccess = 0
}

func (c *Count) Reset() {
	c.RequestFailure, c.RequestSuccess = 0, 0
	c.ConsecutiveFailure, c.ConsecutiveSuccess = 0, 0
}

// Setting defines the configuration parameters for the circuit breaker.
// It includes timing intervals, thresholds for failures and successes, and execution timeouts.
//
// Fields:
//   - Interval: The duration used to check the last failure time. If the time since the last failure exceeds this interval, the state is changed and the failure count is reset.
//   - Timeout: The maximum duration allowed for an execution before timing out.
//   - FailureThreshold: The maximum number of tolerated failures before the circuit breaker transitions to an open state.
//   - SuccessRateThreshold: The minimum success rate required to transition the circuit breaker from half-open to closed state.
type Setting struct {
	Interval             time.Duration
	Timeout              time.Duration
	FailureThreshold     int
	SuccessRateThreshold float32
}

func NewSetting() *Setting {
	return &Setting{
		Interval:             defaultInterval,
		Timeout:              defaultTimeout,
		FailureThreshold:     defaultFailureThreshold,
		SuccessRateThreshold: defaultSuccessRateThreshold,
	}
}

func (s *Setting) SetInterval(interval time.Duration) *Setting {
	s.Interval = interval
	return s
}

func (s *Setting) SetTimeout(timeout time.Duration) *Setting {
	s.Timeout = timeout
	return s
}

func (s *Setting) SetFailureThreshold(maxFailure int) *Setting {
	s.FailureThreshold = maxFailure
	return s
}

func (s *Setting) SetSuccessRate(rate float32) *Setting {
	if rate > 1 {
		rate = 1
	} else if rate < 0 {
		rate = defaultSuccessRateThreshold
	}
	s.SuccessRateThreshold = rate
	return s
}

type CircuitBreaker struct {
	mu          sync.Mutex
	state       State
	count       Count
	setting     *Setting
	lastFailure time.Time
}

func NewCircuitBreaker(setting *Setting) *CircuitBreaker {
	return &CircuitBreaker{
		state:   StateClosed,
		setting: setting,
	}
}

func (cb *CircuitBreaker) CurrentState() string {
	return cb.state.String()
}

// calculateSuccessRate computes and returns the success rate of requests as a float32 value.
// The success rate is calculated as the ratio of successful requests to the total number of requests (successes + failures).
// If there have been no requests, it returns 0 to avoid division by zero.
// This method is thread-safe and acquires a lock to ensure consistent access to the request counters.
func (cb *CircuitBreaker) calculateSuccessRate() float32 {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	totalRequest := cb.count.RequestSuccess + cb.count.RequestFailure
	if totalRequest == 0 {
		return 0
	}
	return float32(cb.count.RequestSuccess) / float32(totalRequest)
}

// closedStateExecution executes the provided handler function while the circuit breaker is in the closed state.
// It manages the success and failure counts, updates the last failure time, and transitions the circuit breaker
// to the open state if the failure threshold is reached. If the handler executes successfully, it resets the
// failure count after a specified interval. The function returns the handler's result or an error if execution fails.
func (cb *CircuitBreaker) closedStateExecution(handler HandlerFunc) (any, error) {
	result, err := cb.executionWithTimeout(handler)
	if err != nil {
		cb.mu.Lock()
		cb.count.OnFailure()
		cb.lastFailure = time.Now()
		if cb.count.ConsecutiveFailure >= cb.setting.FailureThreshold {
			cb.state = StateOpen
		}
		cb.mu.Unlock()
		return nil, err
	}

	cb.mu.Lock()
	cb.count.OnSuccess()
	if time.Since(cb.lastFailure) >= cb.setting.Interval {
		cb.count.Reset()
	}
	cb.mu.Unlock()

	return result, nil
}

// openStateExecution handles the execution logic when the circuit breaker is in the open state.
// If the configured interval has passed since the last failure, it resets the failure count,
// transitions the circuit breaker to the half-open state, and allows the request to proceed.
// Otherwise, it blocks the request and returns an error indicating that the circuit is open.
func (cb *CircuitBreaker) openStateExecution() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if time.Since(cb.lastFailure) >= cb.setting.Interval {
		cb.count.Reset()
		cb.state = StateHalfOpen
		return nil
	}

	return errors.New("request blocked due to circuit open")
}

// halfOpenStateExecution attempts to execute the provided handler function while the circuit breaker is in the half-open state.
// If the execution fails, it records the failure, updates the last failure timestamp, and returns the error.
// On success, it updates the success count and evaluates whether the circuit breaker should transition to the closed state
// based on the elapsed time since the last failure or the current success rate threshold.
// Returns the handler's result and any error encountered during execution.
func (cb *CircuitBreaker) halfOpenStateExecution(handler HandlerFunc) (any, error) {
	result, err := cb.executionWithTimeout(handler)
	if err != nil {
		cb.mu.Lock()
		cb.count.OnFailure()
		cb.lastFailure = time.Now()
		cb.mu.Unlock()

		return nil, err
	}

	successRate := cb.calculateSuccessRate()

	cb.mu.Lock()
	cb.count.OnSuccess()
	sinceLastFailure := time.Since(cb.lastFailure)
	if sinceLastFailure >= cb.setting.Interval ||
		successRate >= cb.setting.SuccessRateThreshold {
		cb.state = StateClosed
		cb.count.Reset()
	}
	cb.mu.Unlock()

	return result, nil
}

// Fire executes the provided handler function based on the current state of the circuit breaker.
// It delegates execution to the appropriate state handler (closed, half-open, or open).
// Returns the result of the handler or an error if the circuit is open or the handler fails.
func (cb *CircuitBreaker) Fire(handler HandlerFunc) (any, error) {
	switch cb.state {
	case StateClosed:
		return cb.closedStateExecution(handler)
	case StateHalfOpen:
		return cb.halfOpenStateExecution(handler)
	case StateOpen:
		return nil, cb.openStateExecution()
	}

	return nil, nil
}
