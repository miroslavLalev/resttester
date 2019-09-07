package stress

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

var genReg = regexp.MustCompile("(\\w+)\\[(\\d+),(\\d+)\\]")

// Generator is interface for sequence of numbers.
type Generator interface {
	Next() int
	Reset()
}

// NewGenerator creates new Generator for the given function string.
// The valid options are:
//  - exp[a,b] -> Creates exponential generator of type a*b^x
//  - lin[a,b] -> Creates linear generator of type a*x+b
func NewGenerator(str string) (Generator, error) {
	matches := genReg.FindAllStringSubmatch(str, -1)
	if len(matches) != 1 {
		return nil, fmt.Errorf("Invalid strategy type")
	}

	a, _ := strconv.Atoi(matches[0][2]) // regex guarantees
	b, _ := strconv.Atoi(matches[0][3]) // regex guarantees
	switch matches[0][1] {
	case "exp":
		return NewExponentialGenerator(a, b), nil
	case "lin":
		return NewLinearGenerator(a, b), nil
	default:
		return nil, fmt.Errorf("Invalid strategy type")
	}
}

// ExponentialGenerator generates exponentially increasing sequence of numbers.
type ExponentialGenerator struct {
	base    int
	step    int
	current int
}

// NewExponentialGenerator creates new exponential generator.
func NewExponentialGenerator(base, step int) Generator {
	return &ExponentialGenerator{base: base, step: step}
}

// Next implements Generator.Next
func (eg *ExponentialGenerator) Next() int {
	res := eg.base * int(math.Pow(float64(eg.step), float64(eg.current)))
	eg.current++
	return res
}

// Reset implements Generator.Reset
func (eg *ExponentialGenerator) Reset() {
	eg.current = 0
}

// LinearGenerator generates linearly increasing sequence of numbers.
type LinearGenerator struct {
	a       int
	b       int
	current int
}

// NewLinearGenerator creates new linear generator.
func NewLinearGenerator(a, b int) Generator {
	return &LinearGenerator{a: a, b: b}
}

// Next implements Generator.Next
func (lg *LinearGenerator) Next() int {
	res := lg.a*lg.current + lg.b
	lg.current++
	return res
}

// Reset implements Generator.Reset
func (lg *LinearGenerator) Reset() {
	lg.current = 0
}
