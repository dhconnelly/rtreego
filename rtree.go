// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A library for efficiently storing and querying spatial data.
package rtreego

// Rtree represents an R-tree, a balanced search tree for storing and querying
// spatial objects.  Dim specifies the number of spatial dimensions and
// MinChildren/MaxChildren specify the minimum/maximum branching factors.
type Rtree struct {
	Dim         uint
	MinChildren uint
	MaxChildren uint
	root        *node
	size        int
}

// NewTree creates a new R-tree instance.  
func NewTree(Dim, MinChildren, MaxChildren uint) *Rtree {
	rt := Rtree{Dim, MinChildren, MaxChildren, new(node), 0}
	rt.root.objects = make([]*Spatial, MinChildren)
	return &rt
}

// Size returns the number of objects currently stored in tree.
func (tree *Rtree) Size() int {
	return tree.size
}

// node represents one entry in an R-tree.
type node struct {
	parent  *node
	entries []entry // non-nil if this is an internal node
	objects []*Spatial // non-nil if this is a leaf node
}

// entry represents one entry in an R-tree.
type entry struct {
	bb     Rect     // bounding-box of all children of this entry
	child  *node
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
	return nil
}

// adjustTree splits overflowing nodes and propagates the changes downwards.
func (tree *Rtree) adjustTree(n *node) {

}

// split splits an overflowing node into two nodes while attempting to minimize
// the area of the resulting nodes.
func (n *node) split() (left, right *node) {
	return nil, nil
}

// pickSeeds chooses the two child nodes of n to start a split.
func (n *node) pickSeeds() (left, right *node) {
	return nil, nil
}

// pickNext chooses a child node of n to be added to left or right.
func (n *node) pickNext(left, right *node) *node {
	return nil
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
