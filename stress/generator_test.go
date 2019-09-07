package stress_test

import (
	"math"
	"resttester/stress"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		str string
		err *string
	}{
		{
			str: "lin[1,2]",
		},
		{
			str: "lin[1, 2]",
			err: StringPtr("Invalid strategy type"),
		},
		{
			str: "lin[1,2,3]",
			err: StringPtr("Invalid strategy type"),
		},
		{
			str: "exp[1,2]",
		},
		{
			str: "wrong[1,2]",
			err: StringPtr("Invalid strategy type"),
		},
	}

	for _, tc := range tests {
		_, err := stress.NewGenerator(tc.str)
		if tc.err == nil {
			AssertNil(t, err)
			continue
		}
		AssertEqual(t, err.Error(), *tc.err)
	}
}

func TestLinearGenerator(t *testing.T) {
	gen, err := stress.NewGenerator("lin[1,1]")
	AssertNil(t, err)

	for i := 0; i < 100; i++ {
		AssertEqual(t, gen.Next(), i+1)
	}
}

func TestExponentialGenerator(t *testing.T) {
	gen, err := stress.NewGenerator("exp[2,2]")
	AssertNil(t, err)

	for i := 0; i < 10; i++ {
		AssertEqual(t, gen.Next(), 2*int(math.Pow(2, float64(i))))
	}
}
