package rtreego

import (
	"fmt"
	"math"
)

// Point represents a point in n-dimensional Euclidean space.
type Point []float64

// Rect represents a subset of n-dimensional Euclidean space of the form
// [a1, b1] x [a2, b2] x ... x [an, bn], where ai < bi for all 1 <= i <= n.
type Rect struct {
	p, q Point // p[i] <= q[i] for all i
}

// DimError represents a failure due to mismatched dimensions.
type DimError struct {
	Expected int
	Actual int
}

func (err *DimError) Error() string {
	return fmt.Sprintf("Expected dim %d, got %d", err.Expected, err.Actual)
}

// Dist computes the Euclidean distance between two points p and q.  When
// len(p) != len(q), a pointer to a DimError instance is returned.
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

// NewRect constructs and returns a pointer to a Rect given a corner point and
// the lengths of each dimension.  The point p should be the most-negative point
// on the rectangle (in every dimension).  If len(lengths) != len(p), then a
// pointer to a DimError instance is returned.
func NewRect(p Point, lengths []float64) (*Rect, error) {
	if len(p) != len(lengths) {
		return nil, &DimError{len(p), len(lengths)}
	}
	q := make([]float64, len(p))
	for i := range p {
		q[i] = p[i] + lengths[i]
	}
	return &Rect{p, q}, nil
}
