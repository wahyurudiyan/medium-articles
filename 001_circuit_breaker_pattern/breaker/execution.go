package breaker

import "context"

type data struct {
	result any
	err    error
}

func (b *CircuitBreaker) executionWithTimeout(handler HandlerFunc) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.setting.Timeout)
	defer cancel()

	execResult := make(chan data, 1)
	go func(ctx context.Context) {
		result, err := handler(ctx)
		execResult <- data{
			result: result, err: err,
		}
	}(ctx)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case data := <-execResult:
		return data.result, data.err
	}
}
