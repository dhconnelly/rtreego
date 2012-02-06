package rtreego

import (
    "math"
    "fmt"
)

type Point struct {
    X float64
    Y float64
}

func NewPoint(x, y float64) *Point {
    return &Point{x, y}
}

func Dist(p, q *Point) float64 {
    dx := p.X - q.X
    dy := p.Y - q.Y
    return math.Sqrt(dx*dx + dy*dy)
}

func (p *Point) String() string {
    return fmt.Sprintf("Point(%f, %f)", p.X, p.Y)
}
