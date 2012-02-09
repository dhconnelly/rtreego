package rtreego

import (
	"math"
	"testing"
)

func TestDist(t *testing.T) {
	p := Point{1, 2, 3}
	q := Point{4, 5, 6}
	r := Point{7, 8}
	dist := math.Sqrt(27)

	if d, err := Dist(p, q); err != nil || d != dist {
		t.Errorf("Dist(%v, %v) = %v; expected %v", p, q, d, dist)
	}

	if d, err := Dist(p, r); d != 0 || err == nil {
		t.Errorf("Expected failure on Dist(%v, %v); got %v", p, r, d)
	}
}

func TestNewRect(t *testing.T) {
	p := Point{1.0, -2.5, 3.0}
	q := Point{3.5, 5.5, 4.5}
	lengths := []float64{2.5, 8.0, 1.5}

	rect, err := NewRect(p, lengths)
	if err != nil {
		t.Errorf("Unexpected failure on NewRect(%v, %v)", p, lengths)
	}
	
	for i, p1 := range p {
		p2 := rect.p[i]
		if p1 != p2 {
			t.Errorf("Expected rect.s[i] = %v; got %v", p1, p2)
		}

		q1 := q[i]
		q2 := rect.q[i]
		if q1 != q2 {
			t.Errorf("Expected rect.s[i] = %v; got %v", q1, q2)
		}
	}

	r := Point{-7.0, 10.0}
	if rect2, err := NewRect(r, lengths); rect2 != nil || err == nil {
		t.Errorf("Expected failure on NewRect(%v, %v); got %v", r, lengths, rect2)
	}
}
