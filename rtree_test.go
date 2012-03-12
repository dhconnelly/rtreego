package rtreego

import (
	"testing"
)

func (r *Rect) Bounds() *Rect {
	return r
}

func TestChooseLeafEmpty(t *testing.T) {
	rt := NewTree(3, 5, 10)

	p := Point{1, 1, 1}
	r, _ := NewRect(p, []float64{1, 1, 1})

	
	if leaf := rt.chooseLeaf(rt.root, r); leaf != rt.root {
		t.Errorf("expected chooseLeaf of empty tree to return root")
	}
}
