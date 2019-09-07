package stress_test

import (
	"resttester/stress"
	"sync/atomic"
	"testing"
)

func TestStrategyExecutions(t *testing.T) {
	s, err := stress.NewStrategy("lin[2,1]")
	AssertNil(t, err)

	for i := 0; i < 50; i++ {
		var n int32

		errs := s.ExecBatch(func() error {
			atomic.AddInt32(&n, 1)
			return nil
		})
		AssertNil(t, errs)
		AssertEqual(t, n, int32(2*i+1))
	}
}
