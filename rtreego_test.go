package rtreego

import (
    "testing"
)

type distTest struct {
    p1, p2 *Point
    out float64
}

var distTests = []distTest {
    distTest{NewPoint(0, 0), NewPoint(3, 4), 5},
}

func TestDist(t *testing.T) {
    for _, dt := range distTests {
        d := Dist(dt.p1, dt.p2)
        if d != dt.out {
            t.Errorf("Dist(%s, %s) = %f, expected %f.", dt.p1, dt.p2, d, dt.out)
        }
    }
}
