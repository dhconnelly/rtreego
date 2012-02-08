package rtreego

import (
	"math"
	"testing"
)

type distTest struct {
	p1, p2 Point
	out float64
	err error
}

var distTests = []distTest{
	distTest{Point{1, 2, 3}, Point{3, 4, 5}, math.Sqrt(12.0), nil},
	distTest{Point{1, 2, 3}, Point{3, 4}, 0, &DimError{3, 2}},
	distTest{Point{1, 2}, Point{3, 4, 5}, 0, &DimError{2, 3}},
}

func TestDist(t *testing.T) {
	for _, dt := range distTests {
		d, err := Dist(dt.p1, dt.p2)
		if d != dt.out {
			t.Errorf("Expected %v, got %v", d, dt.out)
		}

		e1, ok1 := dt.err.(*DimError)
		e2, ok2 := err.(*DimError)
		if (ok1 != ok2) || (ok1 && *e1 != *e2) {
			t.Errorf("Expected '%v', got '%v'", e1, e2)
		}
	}
}
