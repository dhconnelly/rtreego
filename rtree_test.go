package rtreego

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type testCase struct {
	name string
	tree *Rtree
}

func tests(dim, min, max int, objs ...Spatial) []*testCase {
	return []*testCase{
		{
			"dynamically built",
			func() *Rtree {
				rt := NewTree(dim, min, max)
				for _, thing := range objs {
					rt.Insert(thing)
				}
				return rt
			}(),
		},
		{
			"bulk-loaded",
			func() *Rtree {
				return NewTree(dim, min, max, objs...)
			}(),
		},
	}
}

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

func printNode(n *node, level int) {
	padding := strings.Repeat("\t", level)
	fmt.Printf("%sNode: %p\n", padding, n)
	fmt.Printf("%sParent: %p\n", padding, n.parent)
	fmt.Printf("%sLevel: %d\n", padding, n.level)
	fmt.Printf("%sLeaf: %t\n%sEntries:\n", padding, n.leaf, padding)
	for _, e := range n.entries {
		printEntry(e, level+1)
	}
}

func printEntry(e entry, level int) {
	padding := strings.Repeat("\t", level)
	fmt.Printf("%sBB: %v\n", padding, e.bb)
	if e.child != nil {
		printNode(e.child, level)
	} else {
		fmt.Printf("%sObject: %v\n", padding, e.obj)
	}
	fmt.Println()
}

func items(n *node) chan Spatial {
	ch := make(chan Spatial)
	go func() {
		for _, e := range n.entries {
			if n.leaf {
				ch <- e.obj
			} else {
				for obj := range items(e.child) {
					ch <- obj
				}
			}
		}
		close(ch)
	}()
	return ch
}

func verify(t *testing.T, n *node) {
	if n.leaf {
		return
	}
	for _, e := range n.entries {
		if e.child.level != n.level-1 {
			t.Errorf("failed to preserve level order")
		}
		if e.child.parent != n {
			t.Errorf("failed to update parent pointer")
		}
		verify(t, e.child)
	}
}

func indexOf(objs []Spatial, obj Spatial) int {
	ind := -1
	for i, r := range objs {
		if r == obj {
			ind = i
			break
		}
	}
	return ind
}

var chooseLeafNodeTests = []struct {
	bb0, bb1, bb2 *Rect // leaf bounding boxes
	exp           int   // expected chosen leaf
	desc          string
	level         int
}{
	{
		mustRect(Point{1, 1, 1}, []float64{1, 1, 1}),
		mustRect(Point{-1, -1, -1}, []float64{0.5, 0.5, 0.5}),
		mustRect(Point{3, 4, -5}, []float64{2, 0.9, 8}),
		1,
		"clear winner",
		1,
	},
	{
		mustRect(Point{-1, -1.5, -1}, []float64{0.5, 2.5025, 0.5}),
		mustRect(Point{0.5, 1, 0.5}, []float64{0.5, 0.815, 0.5}),
		mustRect(Point{3, 4, -5}, []float64{2, 0.9, 8}),
		1,
		"leaves tie",
		1,
	},
	{
		mustRect(Point{-1, -1.5, -1}, []float64{0.5, 2.5025, 0.5}),
		mustRect(Point{0.5, 1, 0.5}, []float64{0.5, 0.815, 0.5}),
		mustRect(Point{-1, -2, -3}, []float64{2, 4, 6}),
		2,
		"leaf contains obj",
		1,
	},
}

func TestChooseLeafNodeEmpty(t *testing.T) {
	rt := NewTree(3, 5, 10)
	obj := Point{0, 0, 0}.ToRect(0.5)
	e := entry{obj, nil, obj}
	if leaf := rt.chooseNode(rt.root, e, 1); leaf != rt.root {
		t.Errorf("expected chooseLeaf of empty tree to return root")
	}
}

func TestChooseLeafNode(t *testing.T) {
	for _, test := range chooseLeafNodeTests {
		rt := Rtree{}
		rt.root = &node{}

		leaf0 := &node{rt.root, true, []entry{}, 1}
		entry0 := entry{test.bb0, leaf0, nil}

		leaf1 := &node{rt.root, true, []entry{}, 1}
		entry1 := entry{test.bb1, leaf1, nil}

		leaf2 := &node{rt.root, true, []entry{}, 1}
		entry2 := entry{test.bb2, leaf2, nil}

		rt.root.entries = []entry{entry0, entry1, entry2}

		obj := Point{0, 0, 0}.ToRect(0.5)
		e := entry{obj, nil, obj}

		expected := rt.root.entries[test.exp].child
		if leaf := rt.chooseNode(rt.root, e, 1); leaf != expected {
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
	leftEntry := entry{bb: mustRect(Point{1, 1}, []float64{1, 1})}
	left := &node{entries: []entry{leftEntry}}

	rightEntry := entry{bb: mustRect(Point{-1, -1}, []float64{1, 2})}
	right := &node{entries: []entry{rightEntry}}

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

	lbb := l.computeBoundingBox()
	rbb := r.computeBoundingBox()
	if lbb.p.dist(expLeft.p) >= EPS || lbb.q.dist(expLeft.q) >= EPS {
		t.Errorf("expected left.bb = %s, got %s", expLeft, lbb)
	}
	if rbb.p.dist(expRight.p) >= EPS || rbb.q.dist(expRight.q) >= EPS {
		t.Errorf("expected right.bb = %s, got %s", expRight, rbb)
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

	if len(l.entries) != 3 || len(r.entries) != 2 {
		t.Errorf("expected underflow assignment for right group")
	}
}

func TestAssignGroupLeastEnlargement(t *testing.T) {
	r00 := entry{bb: mustRect(Point{0, 0}, []float64{1, 1})}
	r01 := entry{bb: mustRect(Point{0, 1}, []float64{1, 1})}
	r10 := entry{bb: mustRect(Point{1, 0}, []float64{1, 1})}
	r11 := entry{bb: mustRect(Point{1, 1}, []float64{1, 1})}
	r02 := entry{bb: mustRect(Point{0, 2}, []float64{1, 1})}

	group1 := &node{entries: []entry{r00, r01}}
	group2 := &node{entries: []entry{r10, r11}}

	assignGroup(r02, group1, group2)
	if len(group1.entries) != 3 || len(group2.entries) != 2 {
		t.Errorf("expected r02 added to group 1")
	}
}

func TestAssignGroupSmallerArea(t *testing.T) {
	r00 := entry{bb: mustRect(Point{0, 0}, []float64{1, 1})}
	r01 := entry{bb: mustRect(Point{0, 1}, []float64{1, 1})}
	r12 := entry{bb: mustRect(Point{1, 2}, []float64{1, 1})}
	r02 := entry{bb: mustRect(Point{0, 2}, []float64{1, 1})}

	group1 := &node{entries: []entry{r00, r01}}
	group2 := &node{entries: []entry{r12}}

	assignGroup(r02, group1, group2)
	if len(group2.entries) != 2 || len(group1.entries) != 2 {
		t.Errorf("expected r02 added to group 2")
	}
}

func TestAssignGroupFewerEntries(t *testing.T) {
	r0001 := entry{bb: mustRect(Point{0, 0}, []float64{1, 2})}
	r12 := entry{bb: mustRect(Point{1, 2}, []float64{1, 1})}
	r22 := entry{bb: mustRect(Point{2, 2}, []float64{1, 1})}
	r02 := entry{bb: mustRect(Point{0, 2}, []float64{1, 1})}

	group1 := &node{entries: []entry{r0001}}
	group2 := &node{entries: []entry{r12, r22}}

	assignGroup(r02, group1, group2)
	if len(group2.entries) != 2 || len(group1.entries) != 2 {
		t.Errorf("expected r02 added to group 2")
	}
}

func TestAdjustTreeNoPreviousSplit(t *testing.T) {
	rt := Rtree{root: &node{}}

	r00 := entry{bb: mustRect(Point{0, 0}, []float64{1, 1})}
	r01 := entry{bb: mustRect(Point{0, 1}, []float64{1, 1})}
	r10 := entry{bb: mustRect(Point{1, 0}, []float64{1, 1})}
	entries := []entry{r00, r01, r10}
	n := node{rt.root, false, entries, 1}
	rt.root.entries = []entry{entry{bb: Point{0, 0}.ToRect(0), child: &n}}

	rt.adjustTree(&n, nil)

	e := rt.root.entries[0]
	p, q := Point{0, 0}, Point{2, 2}
	if p.dist(e.bb.p) >= EPS || q.dist(e.bb.q) >= EPS {
		t.Errorf("Expected adjustTree to fit %v,%v,%v", r00.bb, r01.bb, r10.bb)
	}
}

func TestAdjustTreeNoSplit(t *testing.T) {
	rt := NewTree(2, 3, 3)

	r00 := entry{bb: mustRect(Point{0, 0}, []float64{1, 1})}
	r01 := entry{bb: mustRect(Point{0, 1}, []float64{1, 1})}
	left := node{rt.root, false, []entry{r00, r01}, 1}
	leftEntry := entry{bb: Point{0, 0}.ToRect(0), child: &left}

	r10 := entry{bb: mustRect(Point{1, 0}, []float64{1, 1})}
	r11 := entry{bb: mustRect(Point{1, 1}, []float64{1, 1})}
	right := node{rt.root, false, []entry{r10, r11}, 1}

	rt.root.entries = []entry{leftEntry}
	retl, retr := rt.adjustTree(&left, &right)

	if retl != rt.root || retr != nil {
		t.Errorf("Expected adjustTree didn't split the root")
	}

	entries := rt.root.entries
	if entries[0].child != &left || entries[1].child != &right {
		t.Errorf("Expected adjustTree keeps left and adds n in parent")
	}

	lbb, rbb := entries[0].bb, entries[1].bb
	if lbb.p.dist(Point{0, 0}) >= EPS || lbb.q.dist(Point{1, 2}) >= EPS {
		t.Errorf("Expected adjustTree to adjust left bb")
	}
	if rbb.p.dist(Point{1, 0}) >= EPS || rbb.q.dist(Point{2, 2}) >= EPS {
		t.Errorf("Expected adjustTree to adjust right bb")
	}
}

func TestAdjustTreeSplitParent(t *testing.T) {
	rt := NewTree(2, 1, 1)

	r00 := entry{bb: mustRect(Point{0, 0}, []float64{1, 1})}
	r01 := entry{bb: mustRect(Point{0, 1}, []float64{1, 1})}
	left := node{rt.root, false, []entry{r00, r01}, 1}
	leftEntry := entry{bb: Point{0, 0}.ToRect(0), child: &left}

	r10 := entry{bb: mustRect(Point{1, 0}, []float64{1, 1})}
	r11 := entry{bb: mustRect(Point{1, 1}, []float64{1, 1})}
	right := node{rt.root, false, []entry{r10, r11}, 1}

	rt.root.entries = []entry{leftEntry}
	retl, retr := rt.adjustTree(&left, &right)

	if len(retl.entries) != 1 || len(retr.entries) != 1 {
		t.Errorf("Expected adjustTree distributed the entries")
	}

	lbb, rbb := retl.entries[0].bb, retr.entries[0].bb
	if lbb.p.dist(Point{0, 0}) >= EPS || lbb.q.dist(Point{1, 2}) >= EPS {
		t.Errorf("Expected left split got left entry")
	}
	if rbb.p.dist(Point{1, 0}) >= EPS || rbb.q.dist(Point{2, 2}) >= EPS {
		t.Errorf("Expected right split got right entry")
	}
}

func TestInsertRepeated(t *testing.T) {
	var things []Spatial
	for i := 0; i < 10; i++ {
		things = append(things, mustRect(Point{0, 0}, []float64{2, 1}))
	}

	for _, tc := range tests(2, 3, 5, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree
			rt.Insert(mustRect(Point{0, 0}, []float64{2, 1}))
		})
	}
}

func TestInsertNoSplit(t *testing.T) {
	rt := NewTree(2, 3, 3)
	thing := mustRect(Point{0, 0}, []float64{2, 1})
	rt.Insert(thing)

	if rt.Size() != 1 {
		t.Errorf("Insert failed to increase tree size")
	}

	if len(rt.root.entries) != 1 || rt.root.entries[0].obj.(*Rect) != thing {
		t.Errorf("Insert failed to insert thing into root entries")
	}
}

func TestInsertSplitRoot(t *testing.T) {
	rt := NewTree(2, 3, 3)
	things := []*Rect{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
	}
	for _, thing := range things {
		rt.Insert(thing)
	}

	if rt.Size() != 6 {
		t.Errorf("Insert failed to insert")
	}

	if len(rt.root.entries) != 2 {
		t.Errorf("Insert failed to split")
	}

	left, right := rt.root.entries[0].child, rt.root.entries[1].child
	if len(left.entries) != 3 || len(right.entries) != 3 {
		t.Errorf("Insert failed to split evenly")
	}
}

func TestInsertSplit(t *testing.T) {
	rt := NewTree(2, 3, 3)
	things := []*Rect{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{10, 10}, []float64{2, 2}),
	}
	for _, thing := range things {
		rt.Insert(thing)
	}

	if rt.Size() != 7 {
		t.Errorf("Insert failed to insert")
	}

	if len(rt.root.entries) != 3 {
		t.Errorf("Insert failed to split")
	}

	a, b, c := rt.root.entries[0], rt.root.entries[1], rt.root.entries[2]
	if len(a.child.entries) != 3 ||
		len(b.child.entries) != 3 ||
		len(c.child.entries) != 1 {
		t.Errorf("Insert failed to split evenly")
	}
}

func TestInsertSplitSecondLevel(t *testing.T) {
	rt := NewTree(2, 3, 3)
	things := []*Rect{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{0, 6}, []float64{1, 2}),
		mustRect(Point{1, 6}, []float64{1, 2}),
		mustRect(Point{0, 8}, []float64{1, 2}),
		mustRect(Point{1, 8}, []float64{1, 2}),
	}
	for _, thing := range things {
		rt.Insert(thing)
	}

	if rt.Size() != 10 {
		t.Errorf("Insert failed to insert")
	}

	// should split root
	if len(rt.root.entries) != 2 {
		t.Errorf("Insert failed to split the root")
	}

	// split level + entries level + objs level
	if rt.Depth() != 3 {
		t.Errorf("Insert failed to adjust properly")
	}

	var checkParents func(n *node)
	checkParents = func(n *node) {
		if n.leaf {
			return
		}
		for _, e := range n.entries {
			if e.child.parent != n {
				t.Errorf("Insert failed to update parent pointers")
			}
			checkParents(e.child)
		}
	}
	checkParents(rt.root)
}

func TestFindLeaf(t *testing.T) {
	rt := NewTree(2, 3, 3)
	things := []*Rect{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{0, 6}, []float64{1, 2}),
		mustRect(Point{1, 6}, []float64{1, 2}),
		mustRect(Point{0, 8}, []float64{1, 2}),
		mustRect(Point{1, 8}, []float64{1, 2}),
	}
	for _, thing := range things {
		rt.Insert(thing)
	}
	verify(t, rt.root)
	for _, thing := range things {
		leaf := rt.findLeaf(rt.root, thing, defaultComparator)
		if leaf == nil {
			printNode(rt.root, 0)
			t.Errorf("Unable to find leaf containing an entry after insertion!")
		}
		var found *Rect
		for _, other := range leaf.entries {
			if other.obj == thing {
				found = other.obj.(*Rect)
				break
			}
		}
		if found == nil {
			printNode(rt.root, 0)
			printNode(leaf, 0)
			t.Errorf("Entry %v not found in leaf node %v!", thing, leaf)
		}
	}
}

func TestFindLeafDoesNotExist(t *testing.T) {
	rt := NewTree(2, 3, 3)
	things := []*Rect{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{0, 6}, []float64{1, 2}),
		mustRect(Point{1, 6}, []float64{1, 2}),
		mustRect(Point{0, 8}, []float64{1, 2}),
		mustRect(Point{1, 8}, []float64{1, 2}),
	}
	for _, thing := range things {
		rt.Insert(thing)
	}

	obj := mustRect(Point{99, 99}, []float64{99, 99})
	leaf := rt.findLeaf(rt.root, obj, defaultComparator)
	if leaf != nil {
		t.Errorf("findLeaf failed to return nil for non-existent object")
	}
}

func TestCondenseTreeEliminate(t *testing.T) {
	rt := NewTree(2, 3, 3)
	things := []*Rect{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{0, 6}, []float64{1, 2}),
		mustRect(Point{1, 6}, []float64{1, 2}),
		mustRect(Point{0, 8}, []float64{1, 2}),
		mustRect(Point{1, 8}, []float64{1, 2}),
	}
	for _, thing := range things {
		rt.Insert(thing)
	}

	// delete entry 2 from parent entries
	parent := rt.root.entries[0].child.entries[1].child
	parent.entries = append(parent.entries[:2], parent.entries[3:]...)
	rt.condenseTree(parent)

	retrieved := []Spatial{}
	for obj := range items(rt.root) {
		retrieved = append(retrieved, obj)
	}

	if len(retrieved) != len(things)-1 {
		t.Errorf("condenseTree failed to reinsert upstream elements")
	}

	verify(t, rt.root)
}

func TestChooseNodeNonLeaf(t *testing.T) {
	rt := NewTree(2, 3, 3)
	things := []*Rect{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{0, 6}, []float64{1, 2}),
		mustRect(Point{1, 6}, []float64{1, 2}),
		mustRect(Point{0, 8}, []float64{1, 2}),
		mustRect(Point{1, 8}, []float64{1, 2}),
	}
	for _, thing := range things {
		rt.Insert(thing)
	}

	obj := mustRect(Point{0, 10}, []float64{1, 2})
	e := entry{obj, nil, obj}
	n := rt.chooseNode(rt.root, e, 2)
	if n.level != 2 {
		t.Errorf("chooseNode failed to stop at desired level")
	}
}

func TestInsertNonLeaf(t *testing.T) {
	rt := NewTree(2, 3, 3)
	things := []*Rect{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{0, 6}, []float64{1, 2}),
		mustRect(Point{1, 6}, []float64{1, 2}),
		mustRect(Point{0, 8}, []float64{1, 2}),
		mustRect(Point{1, 8}, []float64{1, 2}),
	}
	for _, thing := range things {
		rt.Insert(thing)
	}

	obj := mustRect(Point{99, 99}, []float64{99, 99})
	e := entry{obj, nil, obj}
	rt.insert(e, 2)

	expected := rt.root.entries[1].child
	if expected.entries[1].obj != obj {
		t.Errorf("insert failed to insert entry at correct level")
	}
}

func TestDeleteFlatten(t *testing.T) {
	things := []Spatial{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
	}

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree
			// make sure flattening didn't nuke the tree
			rt.Delete(things[0])
			verify(t, rt.root)
		})
	}
}

func TestDelete(t *testing.T) {
	things := []Spatial{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{0, 6}, []float64{1, 2}),
		mustRect(Point{1, 6}, []float64{1, 2}),
		mustRect(Point{0, 8}, []float64{1, 2}),
		mustRect(Point{1, 8}, []float64{1, 2}),
	}

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree

			verify(t, rt.root)

			things2 := []Spatial{}
			for len(things) > 0 {
				i := rand.Int() % len(things)
				things2 = append(things2, things[i])
				things = append(things[:i], things[i+1:]...)
			}

			for i, thing := range things2 {
				ok := rt.Delete(thing)
				if !ok {
					t.Errorf("Thing %v was not found in tree during deletion", thing)
					return
				}

				if rt.Size() != len(things2)-i-1 {
					t.Errorf("Delete failed to remove %v", thing)
					return
				}
				verify(t, rt.root)
			}
		})
	}
}

func TestDeleteWithDepthChange(t *testing.T) {
	rt := NewTree(2, 3, 3)
	things := []*Rect{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
	}
	for _, thing := range things {
		rt.Insert(thing)
	}

	// delete last item and condense nodes
	rt.Delete(things[3])

	// rt.height should be 1 otherwise insert increases height to 3
	rt.Insert(things[3])

	// and verify would fail
	verify(t, rt.root)
}

func TestDeleteWithComparator(t *testing.T) {
	type IDRect struct {
		ID string
		*Rect
	}

	things := []Spatial{
		&IDRect{"1", mustRect(Point{0, 0}, []float64{2, 1})},
		&IDRect{"2", mustRect(Point{3, 1}, []float64{1, 2})},
		&IDRect{"3", mustRect(Point{1, 2}, []float64{2, 2})},
		&IDRect{"4", mustRect(Point{8, 6}, []float64{1, 1})},
		&IDRect{"5", mustRect(Point{10, 3}, []float64{1, 2})},
		&IDRect{"6", mustRect(Point{11, 7}, []float64{1, 1})},
		&IDRect{"7", mustRect(Point{0, 6}, []float64{1, 2})},
		&IDRect{"8", mustRect(Point{1, 6}, []float64{1, 2})},
		&IDRect{"9", mustRect(Point{0, 8}, []float64{1, 2})},
		&IDRect{"10", mustRect(Point{1, 8}, []float64{1, 2})},
	}

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree

			verify(t, rt.root)

			cmp := func(obj1, obj2 Spatial) bool {
				idr1 := obj1.(*IDRect)
				idr2 := obj2.(*IDRect)
				return idr1.ID == idr2.ID
			}

			things2 := []*IDRect{}
			for len(things) > 0 {
				i := rand.Int() % len(things)
				// make a deep copy
				copy := &IDRect{things[i].(*IDRect).ID, &(*things[i].(*IDRect).Rect)}
				things2 = append(things2, copy)

				if !cmp(things[i], copy) {
					log.Fatalf("expected copy to be equal to the original, original: %v, copy: %v", things[i], copy)
				}

				things = append(things[:i], things[i+1:]...)
			}

			for i, thing := range things2 {
				ok := rt.DeleteWithComparator(thing, cmp)
				if !ok {
					t.Errorf("Thing %v was not found in tree during deletion", thing)
					return
				}

				if rt.Size() != len(things2)-i-1 {
					t.Errorf("Delete failed to remove %v", thing)
					return
				}
				verify(t, rt.root)
			}
		})
	}
}

func TestSearchIntersect(t *testing.T) {
	things := []Spatial{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{2, 6}, []float64{1, 2}),
		mustRect(Point{3, 6}, []float64{1, 2}),
		mustRect(Point{2, 8}, []float64{1, 2}),
		mustRect(Point{3, 8}, []float64{1, 2}),
	}

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree

			p := Point{2, 1.5}
			bb := mustRect(p, []float64{10, 5.5})
			q := rt.SearchIntersect(bb)

			var expected []Spatial
			for _, i := range []int{1, 2, 3, 4, 6, 7} {
				expected = append(expected, things[i])
			}

			ensureDisorderedSubset(t, q, expected)
		})
	}

}

func TestSearchIntersectWithLimit(t *testing.T) {
	things := []Spatial{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{2, 6}, []float64{1, 2}),
		mustRect(Point{3, 6}, []float64{1, 2}),
		mustRect(Point{2, 8}, []float64{1, 2}),
		mustRect(Point{3, 8}, []float64{1, 2}),
	}

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree

			bb := mustRect(Point{2, 1.5}, []float64{10, 5.5})

			// expected contains all the intersecting things
			var expected []Spatial
			for _, i := range []int{1, 2, 6, 7, 3, 4} {
				expected = append(expected, things[i])
			}

			// Loop through all possible limits k of SearchIntersectWithLimit,
			// and test that the results are as expected.
			for k := -1; k <= len(things); k++ {
				q := rt.SearchIntersectWithLimit(k, bb)

				if k == -1 {
					ensureDisorderedSubset(t, q, expected)
					if len(q) != len(expected) {
						t.Fatalf("length of actual (%v) was different from expected (%v)", len(q), len(expected))
					}
				} else if k == 0 {
					if len(q) != 0 {
						t.Fatalf("length of actual (%v) was different from expected (%v)", len(q), len(expected))
					}
				} else if k <= len(expected) {
					ensureDisorderedSubset(t, q, expected)
					if len(q) != k {
						t.Fatalf("length of actual (%v) was different from expected (%v)", len(q), len(expected))
					}
				} else {
					ensureDisorderedSubset(t, q, expected)
					if len(q) != len(expected) {
						t.Fatalf("length of actual (%v) was different from expected (%v)", len(q), len(expected))
					}
				}
			}
		})
	}
}

func TestSearchIntersectWithTestFilter(t *testing.T) {
	things := []Spatial{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{2, 6}, []float64{1, 2}),
		mustRect(Point{3, 6}, []float64{1, 2}),
		mustRect(Point{2, 8}, []float64{1, 2}),
		mustRect(Point{3, 8}, []float64{1, 2}),
	}

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree

			bb := mustRect(Point{2, 1.5}, []float64{10, 5.5})

			// intersecting indexes are 1, 2, 6, 7, 3, 4
			// rects which we do not filter out
			var expected []Spatial
			for _, i := range []int{1, 6, 4} {
				expected = append(expected, things[i])
			}

			// this test filter will only pick the objects that are in expected
			objects := rt.SearchIntersect(bb, func(results []Spatial, object Spatial) (bool, bool) {
				for _, exp := range expected {
					if exp == object {
						return false, false
					}
				}
				return true, false
			})

			ensureDisorderedSubset(t, objects, expected)
		})
	}
}

func TestSearchIntersectNoResults(t *testing.T) {
	things := []Spatial{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{2, 6}, []float64{1, 2}),
		mustRect(Point{3, 6}, []float64{1, 2}),
		mustRect(Point{2, 8}, []float64{1, 2}),
		mustRect(Point{3, 8}, []float64{1, 2}),
	}

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree

			bb := mustRect(Point{99, 99}, []float64{10, 5.5})
			q := rt.SearchIntersect(bb)
			if len(q) != 0 {
				t.Errorf("SearchIntersect failed to return nil slice on failing query")
			}
		})
	}
}

func TestSortEntries(t *testing.T) {
	objs := []*Rect{
		mustRect(Point{1, 1}, []float64{1, 1}),
		mustRect(Point{2, 2}, []float64{1, 1}),
		mustRect(Point{3, 3}, []float64{1, 1})}
	entries := []entry{
		entry{objs[2], nil, objs[2]},
		entry{objs[1], nil, objs[1]},
		entry{objs[0], nil, objs[0]},
	}
	sorted, dists := sortEntries(Point{0, 0}, entries)
	if sorted[0] != entries[2] || sorted[1] != entries[1] || sorted[2] != entries[0] {
		t.Errorf("sortEntries failed")
	}
	if dists[0] != 2 || dists[1] != 8 || dists[2] != 18 {
		t.Errorf("sortEntries failed to calculate proper distances")
	}
}

func TestNearestNeighbor(t *testing.T) {
	things := []Spatial{
		mustRect(Point{1, 1}, []float64{1, 1}),
		mustRect(Point{1, 3}, []float64{1, 1}),
		mustRect(Point{3, 2}, []float64{1, 1}),
		mustRect(Point{-7, -7}, []float64{1, 1}),
		mustRect(Point{7, 7}, []float64{1, 1}),
		mustRect(Point{10, 2}, []float64{1, 1}),
	}

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree

			obj1 := rt.NearestNeighbor(Point{0.5, 0.5})
			obj2 := rt.NearestNeighbor(Point{1.5, 4.5})
			obj3 := rt.NearestNeighbor(Point{5, 2.5})
			obj4 := rt.NearestNeighbor(Point{3.5, 2.5})

			if obj1 != things[0] || obj2 != things[1] || obj3 != things[2] || obj4 != things[2] {
				t.Errorf("NearestNeighbor failed")
			}
		})
	}
}

func TestGetAllBoundingBoxes(t *testing.T) {
	rt1 := NewTree(2, 3, 3)
	rt2 := NewTree(2, 2, 4)
	rt3 := NewTree(2, 4, 8)
	things := []*Rect{
		mustRect(Point{0, 0}, []float64{2, 1}),
		mustRect(Point{3, 1}, []float64{1, 2}),
		mustRect(Point{1, 2}, []float64{2, 2}),
		mustRect(Point{8, 6}, []float64{1, 1}),
		mustRect(Point{10, 3}, []float64{1, 2}),
		mustRect(Point{11, 7}, []float64{1, 1}),
		mustRect(Point{10, 10}, []float64{2, 2}),
		mustRect(Point{2, 3}, []float64{0.5, 1}),
		mustRect(Point{3, 5}, []float64{1.5, 2}),
		mustRect(Point{7, 14}, []float64{2.5, 2}),
		mustRect(Point{15, 6}, []float64{1, 1}),
		mustRect(Point{4, 3}, []float64{1, 2}),
		mustRect(Point{1, 7}, []float64{1, 1}),
		mustRect(Point{10, 5}, []float64{2, 2}),
	}
	for _, thing := range things {
		rt1.Insert(thing)
	}
	for _, thing := range things {
		rt2.Insert(thing)
	}
	for _, thing := range things {
		rt3.Insert(thing)
	}

	if rt1.Size() != 14 {
		t.Errorf("Insert failed to insert")
	}
	if rt2.Size() != 14 {
		t.Errorf("Insert failed to insert")
	}
	if rt3.Size() != 14 {
		t.Errorf("Insert failed to insert")
	}

	rtbb1 := rt1.GetAllBoundingBoxes()
	rtbb2 := rt2.GetAllBoundingBoxes()
	rtbb3 := rt3.GetAllBoundingBoxes()

	if len(rtbb1) != 13 {
		t.Errorf("Failed bounding box traversal expected 13 got " + strconv.Itoa(len(rtbb1)))
	}
	if len(rtbb2) != 7 {
		t.Errorf("Failed bounding box traversal expected 7 got " + strconv.Itoa(len(rtbb2)))
	}
	if len(rtbb3) != 2 {
		t.Errorf("Failed bounding box traversal expected 2 got " + strconv.Itoa(len(rtbb3)))
	}
}

type byMinDist struct {
	r []Spatial
	p Point
}

func (r byMinDist) Less(i, j int) bool {
	return r.p.minDist(r.r[i].Bounds()) < r.p.minDist(r.r[j].Bounds())
}

func (r byMinDist) Len() int {
	return len(r.r)
}

func (r byMinDist) Swap(i, j int) {
	r.r[i], r.r[j] = r.r[j], r.r[i]
}

func TestNearestNeighborsAll(t *testing.T) {
	things := []Spatial{
		mustRect(Point{1, 1}, []float64{1, 1}),
		mustRect(Point{-7, -7}, []float64{1, 1}),
		mustRect(Point{1, 3}, []float64{1, 1}),
		mustRect(Point{7, 7}, []float64{1, 1}),
		mustRect(Point{10, 2}, []float64{1, 1}),
		mustRect(Point{3, 3}, []float64{1, 1}),
	}

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree

			p := Point{0.5, 0.5}
			sort.Sort(byMinDist{things, p})

			objs := rt.NearestNeighbors(len(things), p)
			for i := range things {
				if objs[i] != things[i] {
					t.Errorf("NearestNeighbors failed at index %d: %v != %v", i, objs[i], things[i])
				}
			}

			objs = rt.NearestNeighbors(len(things)+2, p)
			if len(objs) > len(things) {
				t.Errorf("NearestNeighbors failed: too many elements")
			}
		})
	}
}

func TestNearestNeighborsFilters(t *testing.T) {
	things := []Spatial{
		mustRect(Point{1, 1}, []float64{1, 1}),
		mustRect(Point{-7, -7}, []float64{1, 1}),
		mustRect(Point{1, 3}, []float64{1, 1}),
		mustRect(Point{7, 7}, []float64{1, 1}),
		mustRect(Point{10, 2}, []float64{1, 1}),
		mustRect(Point{3, 3}, []float64{1, 1}),
	}

	expected := []Spatial{things[0], things[2], things[3]}

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree

			p := Point{0.5, 0.5}
			sort.Sort(byMinDist{expected, p})

			objs := rt.NearestNeighbors(len(things), p, func(r []Spatial, obj Spatial) (bool, bool) {
				for _, ex := range expected {
					if ex == obj {
						return false, false
					}
				}

				return true, false
			})

			ensureOrderedSubset(t, objs, expected)
		})
	}
}

func TestNearestNeighborsHalf(t *testing.T) {
	things := []Spatial{
		mustRect(Point{1, 1}, []float64{1, 1}),
		mustRect(Point{-7, -7}, []float64{1, 1}),
		mustRect(Point{1, 3}, []float64{1, 1}),
		mustRect(Point{7, 7}, []float64{1, 1}),
		mustRect(Point{10, 2}, []float64{1, 1}),
		mustRect(Point{3, 3}, []float64{1, 1}),
	}

	p := Point{0.5, 0.5}
	sort.Sort(byMinDist{things, p})

	for _, tc := range tests(2, 3, 3, things...) {
		t.Run(tc.name, func(t *testing.T) {
			rt := tc.tree

			objs := rt.NearestNeighbors(3, p)
			for i := range objs {
				if objs[i] != things[i] {
					t.Errorf("NearestNeighbors failed at index %d: %v != %v", i, objs[i], things[i])
				}
			}

			objs = rt.NearestNeighbors(len(things)+2, p)
			if len(objs) > len(things) {
				t.Errorf("NearestNeighbors failed: too many elements")
			}
		})
	}
}

func ensureOrderedSubset(t *testing.T, actual []Spatial, expected []Spatial) {
	for i := range actual {
		if len(expected)-1 < i || actual[i] != expected[i] {
			t.Fatalf("actual is not an ordered subset of expected")
		}
	}
}

func ensureDisorderedSubset(t *testing.T, actual []Spatial, expected []Spatial) {
	for _, obj := range actual {
		if !contains(obj, expected) {
			t.Fatalf("actual contained an object that was not expected: %+v", obj)
		}
	}
}

func contains(obj Spatial, slice []Spatial) bool {
	for _, s := range slice {
		if s == obj {
			return true
		}
	}

	return false
}
