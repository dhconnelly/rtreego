package rtreego

import (
	"testing"
)

func (r *Rect) Bounds() *Rect {
	return r
}

func mustRect(r *Rect, err error) *Rect {
	if err != nil {
		panic(err)
	}
	return r
}

var chooseLeafTests = []struct{
	bb0, bb1, bb2 *Rect // leaf bounding boxes
	exp int // expected chosen leaf
	desc string
}{
	{
		mustRect(NewRect(Point{1, 1, 1}, []float64{1, 1, 1})),
		mustRect(NewRect(Point{-1, -1, -1}, []float64{0.5, 0.5, 0.5})),
		mustRect(NewRect(Point{3, 4, -5}, []float64{2, 0.9, 8})),
		1,
		"clear winner",
	},
	{
		mustRect(NewRect(Point{-1, -1.5, -1}, []float64{0.5, 2.5025, 0.5})),
		mustRect(NewRect(Point{0.5, 1, 0.5}, []float64{0.5, 0.815, 0.5})),
		mustRect(NewRect(Point{3, 4, -5}, []float64{2, 0.9, 8})),
		1,
		"leaves tie",
	},
	{
		mustRect(NewRect(Point{-1, -1.5, -1}, []float64{0.5, 2.5025, 0.5})),
		mustRect(NewRect(Point{0.5, 1, 0.5}, []float64{0.5, 0.815, 0.5})),
		mustRect(NewRect(Point{-1, -2, -3}, []float64{2, 4, 6})),
		2,
		"leaf contains obj",
	},
}

func TestChooseLeafEmpty(t *testing.T) {
	rt := NewTree(3, 5, 10)
	obj := Point{0, 0, 0}.ToRect(0.5)
	if leaf := rt.chooseLeaf(rt.root, obj); leaf != rt.root {
		t.Errorf("expected chooseLeaf of empty tree to return root")
	}
}

func TestChooseLeaf(t *testing.T) {
	for _, test := range chooseLeafTests {
		rt := Rtree{}
		rt.root = new(node)
		
		leaf0 := &node{rt.root, nil, []*Spatial{}}
		entry0 := &entry{test.bb0, leaf0}
		
		leaf1 := &node{rt.root, nil, []*Spatial{}}
		entry1 := &entry{test.bb1, leaf1}
		
		leaf2 := &node{rt.root, nil, []*Spatial{}}
		entry2 := &entry{test.bb2, leaf2}

		rt.root.entries = []*entry{entry0, entry1, entry2}

		obj := Point{0, 0, 0}.ToRect(0.5)

		expected := rt.root.entries[test.exp].child
		if leaf := rt.chooseLeaf(rt.root, obj); leaf != expected {
			t.Errorf("TestChooseLeaf(%s): expected %d", test.desc, test.exp)
		}
		
	}
}

func TestPickSeeds(t *testing.T) {
	entry1 := &entry{mustRect(NewRect(Point{1, 1}, []float64{1, 1})), nil}
	entry2 := &entry{mustRect(NewRect(Point{1, -1}, []float64{2, 1})), nil}
	entry3 := &entry{mustRect(NewRect(Point{-1, -1}, []float64{1, 2})), nil}
	entries := []*entry{entry1, entry2, entry3}
	left, right := pickSeeds(entries)
	if left != entry1 || right != entry3 {
		t.Errorf("TestPickSeeds: expected entries %d, %d", 1, 3)
	}
}
