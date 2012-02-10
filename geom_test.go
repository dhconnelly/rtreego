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
	if _, ok := err.(*DistError); !ok {
		t.Errorf("Expected DistError on NewRect(%v, %v)", p, lengths)
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

	r := Point{4.5, -1.7, -3.2}
	if yes, err := rect.ContainsPoint(r); yes || err != nil {
		t.Errorf("Expected %v doesn't contain %v", rect, q)
	}
}

func TestContainsPointDimError(t *testing.T) {
	p := Point{3.7, -2.4, 0.0}
	lengths := []float64{6.2, 1.1, 4.9}
	rect, _ := NewRect(p, lengths)

	q := Point{4.5, -1.7}
	_, err := rect.ContainsPoint(q)
	if _, ok := err.(*DimError); !ok {
		t.Errorf("Expected DimError on %v.ContainsPoint(%v)", rect, q)
	}
}
