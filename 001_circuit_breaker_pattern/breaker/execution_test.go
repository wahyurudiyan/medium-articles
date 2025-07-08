package breaker

import (
	"context"
	"errors"
	"testing"
	"time"
)

type actualResult struct {
	result any
	err    error
}

func TestProcessWithTimeout(t *testing.T) {
	tests := []struct {
		name         string
		timeout      time.Duration
		handler      HandlerFunc
		wantResult   any
		wantErr      error
		expectErrVal bool
	}{
		{
			name:    "handler completes before timeout",
			timeout: 100 * time.Millisecond,
			handler: func(ctx context.Context) (any, error) {
				return "success", nil
			},
			wantResult:   "success",
			wantErr:      nil,
			expectErrVal: false,
		},
		{
			name:    "handler exceeds timeout",
			timeout: 50 * time.Millisecond,
			handler: func(ctx context.Context) (any, error) {
				time.Sleep(200 * time.Millisecond)
				return "late", nil
			},
			wantResult:   nil,
			wantErr:      context.DeadlineExceeded,
			expectErrVal: true,
		},
		{
			name:    "handler returns error",
			timeout: 100 * time.Millisecond,
			handler: func(ctx context.Context) (any, error) {
				return nil, errors.New("handler error")
			},
			wantResult:   nil,
			wantErr:      errors.New("handler error"),
			expectErrVal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := NewSetting()
			st.SetTimeout(tt.timeout)

			cb := NewCircuitBreaker(st)
			got, err := cb.executionWithTimeout(tt.handler)

			t.Logf("Test - %s", tt.name)
			if tt.expectErrVal {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				// For context.DeadlineExceeded, use errors.Is
				if tt.wantErr == context.DeadlineExceeded {
					if !errors.Is(err, context.DeadlineExceeded) {
						t.Errorf("expected context.DeadlineExceeded, got %v", err)
					}
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.(string) != tt.wantResult {
					t.Errorf("expected result %v, got %v", tt.wantResult, got)
				}
			}
		})
	}
}
