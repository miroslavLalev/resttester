package stress

import (
	"sync"
)

// HandleFunc is function that defines unit of execution.
type HandleFunc func() error

// Strategy is used to start units in parallel.
type Strategy struct {
	gen Generator
}

// NewStrategy creates new strategy for the given generator string function.
func NewStrategy(str string) (*Strategy, error) {
	gen, err := NewGenerator(str)
	if err != nil {
		return nil, err
	}
	return &Strategy{gen: gen}, nil
}

// ExecBatch executes the next batch of units.
func (s *Strategy) ExecBatch(handle HandleFunc) []error {
	current := s.gen.Next()
	var (
		wg       sync.WaitGroup
		mErr     sync.Mutex
		execErrs []error
	)
	wg.Add(current)
	for i := 0; i < current; i++ {
		go func() {
			defer wg.Done()
			err := handle()
			if err != nil {
				mErr.Lock()
				defer mErr.Unlock()

				execErrs = append(execErrs, err)
			}
		}()
	}
	wg.Wait()

	return execErrs
}

// Reset resets the strategy state.
func (s *Strategy) Reset() {
	s.gen.Reset()
}
