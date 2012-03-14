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
			t.Errorf("%s: expected %d", test.desc, test.exp)
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
		t.Errorf("expected entries %d, %d", 1, 3)
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
		t.Errorf("expected entry %d", 3)
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

	l, r := n.split(0) // left=entry2, right=entry4
	expLeft := mustRect(Point{1, -1}, []float64{2, 4})
	expRight := mustRect(Point{-3, -3}, []float64{3, 4})

	if l.bb.p.dist(expLeft.p) >= EPS || l.bb.q.dist(expLeft.q) >= EPS {
		t.Errorf("expected left.bb = %s, got %s", expLeft, l.bb)
	}
	if r.bb.p.dist(expRight.p) >= EPS || r.bb.q.dist(expRight.q) >= EPS {
		t.Errorf("expected right.bb = %s, got %s", expRight, r.bb)
	}
}

func TestSplitUnderflow(t *testing.T) {
	entry1 := entry{bb: mustRect(Point{0, 0}, []float64{1, 1})}
	entry2 := entry{bb: mustRect(Point{0, 1}, []float64{1, 1})}
	entry3 := entry{bb: mustRect(Point{0, 2}, []float64{1, 1})}
	entry4 := entry{bb: mustRect(Point{0, 3}, []float64{1, 1})}
	entry5 := entry{bb: mustRect(Point{-50, -50}, []float64{1, 1})}
	entries := []entry{entry1, entry2, entry3, entry4, entry5}
	n := &node{entries: entries}

	l, r := n.split(2)

	if len(l.child.entries) != 3 || len(r.child.entries) != 2 {
		t.Errorf("expected underflow assignment for right group")
	}
}

func TestAssignGroupLeastEnlargement(t *testing.T) {
	r00 := entry{bb: mustRect(Point{0, 0}, []float64{1, 1})}
	r01 := entry{bb: mustRect(Point{0, 1}, []float64{1, 1})}
	r10 := entry{bb: mustRect(Point{1, 0}, []float64{1, 1})}
	r11 := entry{bb: mustRect(Point{1, 1}, []float64{1, 1})}
	r02 := entry{bb: mustRect(Point{0, 2}, []float64{1, 1})}

	bb1 := boundingBox(r00.bb, r01.bb)
	group1 := entry{bb: bb1, child: &node{entries: []entry{r00, r01}}}
	
	bb2 := boundingBox(r10.bb, r11.bb)
	group2 := entry{bb: bb2, child: &node{entries: []entry{r10, r11}}}

	assignGroup(&r02, &group1, &group2)
	if len(group1.child.entries) != 3 || len(group2.child.entries) != 2 {
		t.Errorf("expected r02 added to group 1")
	}
}

func TestAssignGroupSmallerArea(t *testing.T) {
	r00 := entry{bb: mustRect(Point{0, 0}, []float64{1, 1})}
	r01 := entry{bb: mustRect(Point{0, 1}, []float64{1, 1})}
	r12 := entry{bb: mustRect(Point{1, 2}, []float64{1, 1})}
	r02 := entry{bb: mustRect(Point{0, 2}, []float64{1, 1})}

	bb1 := boundingBox(r00.bb, r01.bb)
	group1 := entry{bb: bb1, child: &node{entries: []entry{r00, r01}}}
	
	bb2 := r12
	group2 := entry{bb: bb2.bb, child: &node{entries: []entry{r12}}}

	assignGroup(&r02, &group1, &group2)
	if len(group2.child.entries) != 2 || len(group1.child.entries) != 2 {
		t.Errorf("expected r02 added to group 2")
	}
}

func TestAssignGroupFewerEntries(t *testing.T) {
	r0001 := entry{bb: mustRect(Point{0, 0}, []float64{1, 2})}
	r12 := entry{bb: mustRect(Point{1, 2}, []float64{1, 1})}
	r22 := entry{bb: mustRect(Point{2, 2}, []float64{1, 1})}
	r02 := entry{bb: mustRect(Point{0, 2}, []float64{1, 1})}

	bb1 := r0001.bb
	group1 := entry{bb: bb1, child: &node{entries: []entry{r0001}}}
	
	bb2 := boundingBox(r12.bb, r22.bb)
	group2 := entry{bb: bb2, child: &node{entries: []entry{r12, r22}}}

	assignGroup(&r02, &group1, &group2)
	if len(group2.child.entries) != 2 || len(group1.child.entries) != 2 {
		t.Errorf("expected r02 added to group 2")
	}
}
