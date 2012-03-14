// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A library for efficiently storing and querying spatial data.
package rtreego

import (
	"math"
)

// Rtree represents an R-tree, a balanced search tree for storing and querying
// spatial objects.  Dim specifies the number of spatial dimensions and
// MinChildren/MaxChildren specify the minimum/maximum branching factors.
type Rtree struct {
	Dim         uint
	MinChildren uint
	MaxChildren uint
	root        node
	size        int
}

// NewTree creates a new R-tree instance.  
func NewTree(Dim, MinChildren, MaxChildren uint) *Rtree {
	rt := Rtree{Dim: Dim, MinChildren: MinChildren, MaxChildren: MaxChildren}
	rt.root.entries = make([]entry, MinChildren)
	rt.root.leaf = true
	return &rt
}

// Size returns the number of objects currently stored in tree.
func (tree *Rtree) Size() int {
	return tree.size
}

// node represents a tree node of an Rtree.
type node struct {
	parent  *node
	leaf bool
	entries []entry
}

// entry represents a spatial index record stored in a tree node.
type entry struct {
	bb     *Rect     // bounding-box of all children of this entry
	child  *node
	obj *Spatial
}

// Any type that implements Spatial can be stored in an Rtree and queried.
type Spatial interface {
	Bounds() *Rect
}

// Insertion

// Insert inserts a spatial object into the tree.  A DimError is returned if
// the dimensions of the object don't match those of the tree.  If insertion
// causes a leaf node to overflow, the tree is rebalanced automatically.
//
// Implemented per Section 3.2 of "R-trees: A Dynamic Index Structure for
// Spatial Searching" by A. Guttman, Proceedings of ACM SIGMOD, p. 47-57, 1984.
func (tree *Rtree) Insert(obj Spatial) error {
	return nil
}

// chooseLeaf finds the leaf node in which obj should be inserted.
func (tree *Rtree) chooseLeaf(n *node, obj Spatial) *node {
	if n.leaf {
		return n
	}

	// find the entry whose bb needs least enlargement to include obj
	diff := math.MaxFloat64
	var chosen entry
	for _, e := range n.entries {
		bb := boundingBox(e.bb, obj.Bounds())
		d := bb.size() - e.bb.size()
		if d < diff || (d == diff && e.bb.size() < chosen.bb.size()) {
			diff = d
			chosen = e
		}
	}

	return tree.chooseLeaf(chosen.child, obj)
}

// adjustTree splits overflowing nodes and propagates the changes upwards.
func (tree *Rtree) adjustTree(n *node) {
	if n == &tree.root {
		return
	}

	n.resizeBoundingBox()
	tree.adjustTree(n.parent)
}

func (tree *Rtree) adjustTreeSplit(n, nn *node) {
	
}

// resizeBoundingBox adjusts the bounding box of a node to its minimum
// bounding rectangle.
func (n *node) resizeBoundingBox() {
	var ownEntry *entry
	for i := range n.parent.entries {
		if n.parent.entries[i].child == n {
			ownEntry = &n.parent.entries[i]
			break
		}
	}
	childBoxes := []*Rect{}
	for _, e := range n.entries {
		childBoxes = append(childBoxes, e.bb)
	}
	ownEntry.bb = boundingBoxN(childBoxes...)
}

// split splits a node into two groups while attempting to minimize the
// bounding-box area of the resulting groups.
func (n *node) split(minGroupSize int) (left, right entry) {
	l, r := n.pickSeeds()
	leftSeed, rightSeed := n.entries[l], n.entries[r]

	// new nodes can't be leaves, even if n is a leaf
	left = entry{leftSeed.bb, &node{entries: []entry{leftSeed}}, nil}
	right = entry{rightSeed.bb, &node{entries: []entry{rightSeed}}, nil}

	// get the entries to be divided between left and right
	remaining := append(n.entries[:l], n.entries[l+1:r]...)
	remaining = append(remaining, n.entries[r+1:]...)
	
	for len(remaining) > 0 {
		next := pickNext(left, right, remaining)
		e := remaining[next]

		// check for underflow
		if len(remaining) + len(left.child.entries) <= minGroupSize {
			assign(&e, &left)
		} else if len(remaining) + len(right.child.entries) <= minGroupSize {
			assign(&e, &right)
		} else {
			assignGroup(&e, &left, &right)
		}

		remaining = append(remaining[:next], remaining[next+1:]...)
	}
	
	return
}

func assign(e, group *entry) {
	group.child.entries = append(group.child.entries, *e)
	group.bb = boundingBox(group.bb, e.bb)
}

// assignGroup chooses one of two groups to which a node should be added.
func assignGroup(e, left, right *entry) {
	leftEnlarged := boundingBox(left.bb, e.bb)
	rightEnlarged := boundingBox(right.bb, e.bb)

	// first, choose the group that needs the least enlargement
	leftDiff := leftEnlarged.size() - left.bb.size()
	rightDiff := rightEnlarged.size() - right.bb.size()
	if diff := leftDiff - rightDiff; diff < 0 {
		assign(e, left)
		return
	} else if diff > 0 {
		assign(e, right)
		return
	}

	// next, choose the group that has smaller area
	if diff := left.bb.size() - right.bb.size(); diff < 0 {
		assign(e, left)
		return
	} else if diff > 0 {
		assign(e, right)
		return
	}

	// next, choose the group with fewer entries
	if diff := len(left.child.entries) - len(right.child.entries); diff <= 0 {
		assign(e, left)
		return
	}
	assign(e, right)
}

// pickSeeds chooses two child entries of n to start a split.
func (n *node) pickSeeds() (left, right int) {
	maxWastedSpace := -1.0
	for i, e1 := range n.entries {
		for j, e2 := range n.entries[i+1:] {
			d := boundingBox(e1.bb, e2.bb).size() - e1.bb.size() - e2.bb.size()
			if d > maxWastedSpace {
				maxWastedSpace = d
				left, right = i, j+i+1
			}
		}
	}
	return
}

// pickNext chooses an entry to be added to an entry group.
func pickNext(left, right entry, entries []entry) (next int) {
	maxDiff := -1.0
	for i, e := range entries {
		d1 := boundingBox(left.bb, e.bb).size() - left.bb.size()
		d2 := boundingBox(right.bb, e.bb).size() - right.bb.size()
		d := math.Abs(d1 - d2)
		if d > maxDiff {
			maxDiff = d
			next = i
		}
	}
	return
}

// Deletion

// Delete removes an object from the tree.  If the object is not found, ok
// is false; otherwise ok is true.  A DimError is returned if the specified
// object has improper dimensions for the tree.
//
// Implemented per Section 3.3 of "R-trees: A Dynamic Index Structure for
// Spatial Searching" by A. Guttman, Proceedings of ACM SIGMOD, p. 47-57, 1984.
func (tree *Rtree) Delete(obj Spatial) (ok bool, err error) {
	return false, nil
}

// findLeaf finds the leaf node containing obj.
func (tree *Rtree) findLeaf(n *node, obj Spatial) *node {
	return nil
}

// condenseTree deletes underflowing nodes and propagates the changes upwards.
func (tree *Rtree) condenseTree(n *node) *node {
	return nil
}
