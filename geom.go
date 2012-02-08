package rtreego

import (
	"fmt"
	"math"
)

type Point []float64

type DimError struct {
	Expected int
	Actual int
}

func (err *DimError) Error() string {
	return fmt.Sprintf("Expected dim %d, got %d", err.Expected, err.Actual)
}

func Dist(p, q Point) (float64, error) {
	if len(p) != len(q) {
		return 0, &DimError{len(p), len(q)}
	}
	sum := 0.0
	for i := range p {
		dx := p[i] - q[i]
		sum += dx * dx
	}
	return math.Sqrt(sum), nil
}
