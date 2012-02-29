// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtreego

import (
	"math"
	"testing"
)

const EPS = 0.000000001

func TestDist(t *testing.T) {
	p := Point{1, 2, 3}
	q := Point{4, 5, 6}
	dist := math.Sqrt(27)
	if d, err := Dist(p, q); err != nil || d != dist {
		t.Errorf("Dist(%v, %v) = %v; expected %v", p, q, d, dist)
	}
}

func TestDistDimMismatch(t *testing.T) {
	p := Point{1, 2, 3}
	r := Point{7, 8}
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
		t.Errorf("Error on NewRect(%v, %v): %v", p, lengths, err)
	}
	if d, _ := Dist(p, rect.p); d > EPS {
		t.Errorf("Expected p == rect.p")
	}
	if d, _ := Dist(q, rect.q); d > EPS {
		t.Errorf("Expected q == rect.q")
	}
}

func TestNewRectDimMismatch(t *testing.T) {
	p := Point{-7.0, 10.0}
	lengths := []float64{2.5, 8.0, 1.5}
	_, err := NewRect(p, lengths)
	if _, ok := err.(*DimError); !ok {
		t.Errorf("Expected DimError on NewRect(%v, %v)", p, lengths)
	}
}

func TestNewRectDistError(t *testing.T) {
	p := Point{1.0, -2.5, 3.0}
	lengths := []float64{2.5, -8.0, 1.5}
	_, err := NewRect(p, lengths)
	if _, ok := err.(DistError); !ok {
		t.Errorf("Expected DistError on NewRect(%v, %v)", p, lengths)
	}
}

func TestRectSize(t *testing.T) {
	p := Point{1.0, -2.5, 3.0}
	lengths := []float64{2.5, 8.0, 1.5}
	rect, _ := NewRect(p, lengths)
	size := lengths[0] * lengths[1] * lengths[2]
	actual := rect.Size()
	if size != actual {
		t.Errorf("Expected %v.Size() == %v, got %v", rect, size, actual)
	}
}

func TestRectMargin(t *testing.T) {
	p := Point{1.0, -2.5, 3.0}
	lengths := []float64{2.5, 8.0, 1.5}
	rect, _ := NewRect(p, lengths)
	size := 4*2.5 + 4*8.0 + 4*1.5
	actual := rect.Margin()
	if size != actual {
		t.Errorf("Expected %v.Margin() == %v, got %v", rect, size, actual)
	}
}

func TestContainsPoint(t *testing.T) {
	p := Point{3.7, -2.4, 0.0}
	lengths := []float64{6.2, 1.1, 4.9}
	rect, _ := NewRect(p, lengths)

	q := Point{4.5, -1.7, 4.8}
	if yes, err := rect.ContainsPoint(q); !yes || err != nil {
		t.Errorf("Expected %v contains %v", rect, q)
	}
}

func TestDoesNotContainPoint(t *testing.T) {
	p := Point{3.7, -2.4, 0.0}
	lengths := []float64{6.2, 1.1, 4.9}
	rect, _ := NewRect(p, lengths)

	q := Point{4.5, -1.7, -3.2}
	if yes, _ := rect.ContainsPoint(q); yes {
		t.Errorf("Expected %v doesn't contain %v", rect, q)
	}
}

func TestContainsRect(t *testing.T) {
	p := Point{3.7, -2.4, 0.0}
	lengths1 := []float64{6.2, 1.1, 4.9}
	rect1, _ := NewRect(p, lengths1)

	q := Point{4.1, -1.9, 1.0}
	lengths2 := []float64{3.2, 0.6, 3.7}
	rect2, _ := NewRect(q, lengths2)

	if yes, err := rect1.ContainsRect(rect2); !yes || err != nil {
		t.Errorf("Expected %v.ContainsRect(%v", rect1, rect2)
	}
}

func TestDoesNotContainRectOverlaps(t *testing.T) {
	p := Point{3.7, -2.4, 0.0}
	lengths1 := []float64{6.2, 1.1, 4.9}
	rect1, _ := NewRect(p, lengths1)

	q := Point{4.1, -1.9, 1.0}
	lengths2 := []float64{3.2, 1.4, 3.7}
	rect2, _ := NewRect(q, lengths2)

	if yes, _ := rect1.ContainsRect(rect2); yes {
		t.Errorf("Expected %v doesn't contain %v", rect1, rect2)
	}
}

func TestDoesNotContainRectDisjoint(t *testing.T) {
	p := Point{3.7, -2.4, 0.0}
	lengths1 := []float64{6.2, 1.1, 4.9}
	rect1, _ := NewRect(p, lengths1)

	q := Point{1.2, -19.6, -4.0}
	lengths2 := []float64{2.2, 5.9, 0.5}
	rect2, _ := NewRect(q, lengths2)

	if yes, _ := rect1.ContainsRect(rect2); yes {
		t.Errorf("Expected %v doesn't contain %v", rect1, rect2)
	}
}

func TestNoIntersection(t *testing.T) {
	p := Point{1, 2, 3}
	lengths1 := []float64{1, 1, 1}
	rect1, _ := NewRect(p, lengths1)

	q := Point{-1, -2, -3}
	lengths2 := []float64{2.5, 3, 6.5}
	rect2, _ := NewRect(q, lengths2)

	// rect1 and rect2 fail to overlap in just one dimension (second)

	if intersect, _ := Intersect(rect1, rect2); intersect != nil {
		t.Errorf("Expected Intersect(%v, %v) == nil, got %v", rect1, rect2, intersect)
	}
}

func TestNoIntersectionJustTouches(t *testing.T) {
	p := Point{1, 2, 3}
	lengths1 := []float64{1, 1, 1}
	rect1, _ := NewRect(p, lengths1)

	q := Point{-1, -2, -3}
	lengths2 := []float64{2.5, 4, 6.5}
	rect2, _ := NewRect(q, lengths2)

	// rect1 and rect2 fail to overlap in just one dimension (second)

	if intersect, _ := Intersect(rect1, rect2); intersect != nil {
		t.Errorf("Expected Intersect(%v, %v) == nil, got %v", rect1, rect2, intersect)
	}
}

func TestContainmentIntersection(t *testing.T) {
	p := Point{1, 2, 3}
	lengths1 := []float64{1, 1, 1}
	rect1, _ := NewRect(p, lengths1)

	q := Point{1, 2.2, 3.3}
	lengths2 := []float64{0.5, 0.5, 0.5}
	rect2, _ := NewRect(q, lengths2)

	r := Point{1, 2.2, 3.3}
	s := Point{1.5, 2.7, 3.8}

	actual, _ := Intersect(rect1, rect2)
	d1, _ := Dist(r, actual.p)
	d2, _ := Dist(s, actual.q)
	if d1 > EPS || d2 > EPS {
		t.Errorf("Intersect(%v, %v) != %v, %v, got %v", rect1, rect2, r, s, actual)
	}
}

func TestOverlapIntersection(t *testing.T) {
	p := Point{1, 2, 3}
	lengths1 := []float64{1, 2.5, 1}
	rect1, _ := NewRect(p, lengths1)

	q := Point{1, 4, -3}
	lengths2 := []float64{3, 2, 6.5}
	rect2, _ := NewRect(q, lengths2)

	r := Point{1, 4, 3}
	s := Point{2, 4.5, 3.5}

	actual, _ := Intersect(rect1, rect2)
	d1, _ := Dist(r, actual.p)
	d2, _ := Dist(s, actual.q)
	if d1 > EPS || d2 > EPS {
		t.Errorf("Intersect(%v, %v) != %v, %v, got %v", rect1, rect2, r, s, actual)
	}
}

func TestToRect(t *testing.T) {
	x := Point{3.7, -2.4, 0.0}
	tol := 0.05
	rect := x.ToRect(tol)

	p := Point{3.65, -2.45, -0.05}
	q := Point{3.75, -2.35, 0.05}
	d1, _ := Dist(p, rect.p)
	d2, _ := Dist(q, rect.q)
	if d1 > EPS || d2 > EPS {
		t.Errorf("Expected %v.ToRect(%v) == %v, %v, got %v", x, tol, p, q, rect)
	}
}

func TestBoundingBox(t *testing.T) {
	p := Point{3.7, -2.4, 0.0}
	lengths1 := []float64{1, 15, 3}
	rect1, _ := NewRect(p, lengths1)

	q := Point{-6.5, 4.7, 2.5}
	lengths2 := []float64{4, 5, 6}
	rect2, _ := NewRect(q, lengths2)

	r := Point{-6.5, -2.4, 0.0}
	s := Point{4.7, 12.6, 8.5}

	bb, _ := BoundingBox(rect1, rect2)
	d1, _ := Dist(r, bb.p)
	d2, _ := Dist(s, bb.q)
	if d1 > EPS || d2 > EPS {
		t.Errorf("BoundingBox(%v, %v) != %v, %v, got %v", rect1, rect2, r, s, bb)
	}
}

func TestBoundingBoxContains(t *testing.T) {
	p := Point{3.7, -2.4, 0.0}
	lengths1 := []float64{1, 15, 3}
	rect1, _ := NewRect(p, lengths1)

	q := Point{4.0, 0.0, 1.5}
	lengths2 := []float64{0.56, 6.222222, 0.946}
	rect2, _ := NewRect(q, lengths2)

	bb, _ := BoundingBox(rect1, rect2)
	d1, _ := Dist(rect1.p, bb.p)
	d2, _ := Dist(rect1.q, bb.q)
	if d1 > EPS || d2 > EPS {
		t.Errorf("BoundingBox(%v, %v) != %v, got %v", rect1, rect2, rect1, bb)
	}
}

func TestMinDistZero(t *testing.T) {
	p := Point{1, 2, 3}
	r := p.ToRect(1)
	if d, _ := MinDist(p, r); d > EPS {
		t.Errorf("Expected MinDist(%v, %v) == 0, got %v", p, r, d)
	}
}

func TestMinDistPositive(t *testing.T) {
	p := Point{1, 2, 3}
	r := &Rect{Point{-1, -4, 7}, Point{2, -2, 9}}
	expected := float64((-2-2)*(-2-2) + (7-3)*(7-3))
	if d, _ := MinDist(p, r); math.Abs(d - expected) > EPS {
		t.Errorf("Expected MinDist(%v, %v) == %v, got %v", p, r, expected, d)
	}
}

func TestMinMaxDist(t *testing.T) {
	p := Point{-3, -2, -1}
	r := &Rect{Point{0, 0, 0}, Point{1, 2, 3}}

	// furthest points from p on the faces closest to p in each dimension
	q1 := Point{0, 2, 3}
	q2 := Point{1, 0, 3}
	q3 := Point{1, 2, 0}

	// find the closest distance from p to one of these furthest points
	d1, _ := Dist(p, q1)
	d2, _ := Dist(p, q2)
	d3, _ := Dist(p, q3)
	expected := math.Min(d1*d1, math.Min(d2*d2, d3*d3))

	if d, _ := MinMaxDist(p, r); math.Abs(d - expected) > EPS {
		t.Errorf("Expected MinMaxDist(%v, %v) == %v, got %v", p, r, expected, d)
	}
}
