package rtreego

import (
	"fmt"
	"strings"
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
		fmt.Printf("%sObject: %p\n", padding, e.obj)
	}
	fmt.Println()
}

var chooseLeafNodeTests = []struct {
	bb0, bb1, bb2 *Rect // leaf bounding boxes
	exp           int   // expected chosen leaf
	desc          string
	level int
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

	obj := mustRect(Point{1.5, 7}, []float64{0.5, 0.5})
	leaf := rt.findLeaf(rt.root, obj)
	expected := rt.root.entries[0].child.entries[1].child
	if leaf != expected {
		t.Errorf("Failed to locate leaf containing %v", obj)
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

	// TODO verify levels
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
