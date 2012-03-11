// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A library for efficiently storing and querying spatial data.
package rtreego

type Rtree struct {
	Dim uint
	MinChildrenPerNode uint
	MaxChildrenPerNode uint
	root *node
	size int
	height int
}

type Spatial interface {
	Bounds() *Rect
}

func NewTree(Dim, MinChildrenPerNode, MaxChildrenPerNode uint) *Rtree {
	return &Rtree{Dim, MinChildrenPerNode, MaxChildrenPerNode, nil, 0, 0}
}

func (tree *Rtree) Size() int {
	return tree.size
}

func (tree *Rtree) Height() int {
	return tree.height
}

type node struct {
	bb *Rect
	children []*node
	objects []*Spatial
}
