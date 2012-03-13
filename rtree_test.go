package rtreego

import (
	"testing"
)

func (r *Rect) Bounds() *Rect {
	return r
}

func mustRect(p Point, widths []float64) *Rect {
	r, err := NewRect(p, widths)
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
		mustRect(Point{1, 1, 1}, []float64{1, 1, 1}),
		mustRect(Point{-1, -1, -1}, []float64{0.5, 0.5, 0.5}),
		mustRect(Point{3, 4, -5}, []float64{2, 0.9, 8}),
		1,
		"clear winner",
	},
	{
		mustRect(Point{-1, -1.5, -1}, []float64{0.5, 2.5025, 0.5}),
		mustRect(Point{0.5, 1, 0.5}, []float64{0.5, 0.815, 0.5}),
		mustRect(Point{3, 4, -5}, []float64{2, 0.9, 8}),
		1,
		"leaves tie",
	},
	{
		mustRect(Point{-1, -1.5, -1}, []float64{0.5, 2.5025, 0.5}),
		mustRect(Point{0.5, 1, 0.5}, []float64{0.5, 0.815, 0.5}),
		mustRect(Point{-1, -2, -3}, []float64{2, 4, 6}),
		2,
		"leaf contains obj",
	},
}

func TestChooseLeafEmpty(t *testing.T) {
	rt := NewTree(3, 5, 10)
	obj := Point{0, 0, 0}.ToRect(0.5)
	if leaf := rt.chooseLeaf(&rt.root, obj); leaf != &rt.root {
		t.Errorf("expected chooseLeaf of empty tree to return root")
	}
}

func TestChooseLeaf(t *testing.T) {
	for _, test := range chooseLeafTests {
		rt := Rtree{}
		rt.root = node{}
		
		leaf0 := &node{&rt.root, true, []entry{}}
		entry0 := entry{test.bb0, leaf0, nil}
		
		leaf1 := &node{&rt.root, true, []entry{}}
		entry1 := entry{test.bb1, leaf1, nil}
		
		leaf2 := &node{&rt.root, true, []entry{}}
		entry2 := entry{test.bb2, leaf2, nil}

		rt.root.entries = []entry{entry0, entry1, entry2}

		obj := Point{0, 0, 0}.ToRect(0.5)

		expected := rt.root.entries[test.exp].child
		if leaf := rt.chooseLeaf(&rt.root, obj); leaf != expected {
			t.Errorf("TestChooseLeaf(%s): expected %d", test.desc, test.exp)
		}
		
	}
}

func TestPickSeeds(t *testing.T) {
	entry1 := entry{bb: mustRect(Point{1, 1}, []float64{1, 1})}
	entry2 := entry{bb: mustRect(Point{1, -1}, []float64{2, 1})}
	entry3 := entry{bb: mustRect(Point{-1, -1}, []float64{1, 2})}
	n := node{entries: []entry{entry1, entry2, entry3}}
	left, right := n.pickSeeds()
	if n.entries[left] != entry1 || n.entries[right] != entry3 {
		t.Errorf("TestPickSeeds: expected entries %d, %d", 1, 3)
	}
}

func TestPickNext(t *testing.T) {
	left := entry{bb: mustRect(Point{1, 1}, []float64{1, 1})}
	right := entry{bb: mustRect(Point{-1, -1}, []float64{1, 2})}

	entry1 := entry{bb: mustRect(Point{0, 0}, []float64{1, 1})}
	entry2 := entry{bb: mustRect(Point{-2, -2}, []float64{1, 1})}
	entry3 := entry{bb: mustRect(Point{1, 2}, []float64{1, 1})}
	entries := []entry{entry1, entry2, entry3}

	chosen := pickNext(left, right, entries)
	if entries[chosen] != entry2 {
		t.Errorf("TestPickNext: expected entry %d", 3)
	}
}

func TestSplit(t *testing.T) {
	entry1 := entry{bb: mustRect(Point{-3, -1}, []float64{2, 1})}
	entry2 := entry{bb: mustRect(Point{1, 2}, []float64{1, 1})}
	entry3 := entry{bb: mustRect(Point{-1, 0}, []float64{1, 1})}
	entry4 := entry{bb: mustRect(Point{-3, -3}, []float64{1, 1})}
	entry5 := entry{bb: mustRect(Point{1, -1}, []float64{2, 2})}
	entries := []entry{entry1, entry2, entry3, entry4, entry5}
	n := &node{entries: entries}

	left, right := n.split(0) // left=entry2, right=entry4
	leftBB := mustRect(Point{1, -1}, []float64{2, 4})
	rightBB := mustRect(Point{-3, -3}, []float64{3, 4})

	dlp, _ := left.bb.p.Dist(leftBB.p)
	dlq, _ := left.bb.q.Dist(leftBB.q)
	if dlp >= EPS || dlq >= EPS {
		t.Errorf("TestSplit: expected left.bb = %s, got %s", leftBB, left.bb)
	}

	drp, _ := right.bb.p.Dist(rightBB.p)
	drq, _ := right.bb.q.Dist(rightBB.q)
	if drp >= EPS || drq >= EPS {
		t.Errorf("TestSplit: expected right.bb = %s, got %s", rightBB, right.bb)
	}
}
