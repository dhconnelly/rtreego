rtreego
=======

A library for efficiently storing and querying spatial data
in the Go programming language.

About
-----

The R-tree is a popular data structure for efficiently storing and
querying spatial objects; one common use is implementing geospatial
indexes in database management systems.  The variant implemented here,
known as the R*-tree, improves performance and increases storage
utilization.  Both bounding-box queries and k-nearest-neighbor queries
are supported.

R-trees are balanced, so maximum tree height is guaranteed to be
logarithmic in the number of entries; however, good worst-case
performance is not guaranteed.  Instead, a number of rebalancing
heuristics are applied that perform well in practice.  For more
details please refer to the references.

This implementation handles the general N-dimensional case; for a more
efficient implementation for the 3-dimensional case, see [Patrick
Higgins' fork](https://github.com/patrick-higgins/rtreego).

Getting Started
---------------

Get the source code from [GitHub](https://github.com/dhconnelly/rtreego) or,
with Go 1 installed, run `go get github.com/dhconnelly/rtreego`.

Make sure you `import github.com/dhconnelly/rtreego` in your Go source files.

Documentation
-------------

### Storing, updating, and deleting objects

To create a new tree, specify the number of spatial dimensions and the minimum
and maximum branching factor:

	rt := rtreego.NewTree(2, 25, 50)

Any type that implements the `Spatial` interface can be stored in the tree:

	type Spatial interface {
		Bounds() *Rect
	}

`Rect`s are data structures for representing spatial objects, while `Point`s
represent spatial locations.  Creating `Point`s is easy--they're just slices
of `float64`s:

	p1 := rtreego.Point{0.4, 0.5}
	p2 := rtreego.Point{6.2, -3.4}

To create a `Rect`, specify a location and the lengths of the sides:

	r1 := rtreego.NewRect(p1, []float64{1, 2})
	r2 := rtreego.NewRect(p2, []float64{1.7, 2.7})

To demonstrate, let's create and store some test data.

	type Thing struct {
		where *Rect
		name string
	}

	func (t *Thing) Bounds() *Rect {
		return t.where
	}

	rt.Insert(&Thing{r1, "foo"})
	rt.Insert(&Thing{r2, "bar"})

	size := rt.Size() // returns 2

We can insert and delete objects from the tree in any order.

	rt.Delete(thing2)
	// do some stuff...
	rt.Insert(anotherThing)

If you want to store points instead of rectangles, you can easily convert a
point into a rectangle using the `ToRect` method:

	var tol = 0.01

	type Somewhere struct {
		location rtreego.Point
		name string
		wormhole chan int
	}

	func (s *Somewhere) Bounds() *Rect {
		// define the bounds of s to be a rectangle centered at s.location
		// with side lengths 2 * tol:
		return s.location.ToRect(tol)
	}

	rt.Insert(&Somewhere{rtreego.Point{0, 0}, "Someplace", nil})

If you want to update the location of an object, you must delete it, update it,
and re-insert.  Just modifying the object so that the `*Rect` returned by
`Location()` changes, without deleting and re-inserting the object, will
corrupt the tree.

### Queries

Bounding-box and k-nearest-neighbors queries are supported.

Bounding-box queries require a search `*Rect` argument and come in two flavors:
containment search and intersection search.  The former returns all objects that
fall strictly inside the search rectangle, while the latter returns all objects
that touch the search rectangle.

	bb := rtreego.NewRect(rtreego.Point{1.7, -3.4}, []float64{3.2, 1.9})

	// Get a slice of the objects in rt that intersect bb:
	results, _ := rt.SearchIntersect(bb)

	// Get a slice of the objects in rt that are contained inside bb:
	results, _ = rt.SearchContained(bb)

Nearest-neighbor queries find the objects in a tree closest to a specified
query point.

	q := rtreego.Point{6.5, -2.47}
	k := 5

	// Get a slice of the k objects in rt closest to q:
	results, _ = rt.SearchNearestNeighbors(q, k)

### More information

See [GoPkgDoc](http://gopkgdoc.appspot.com/pkg/github.com/dhconnelly/rtreego)
for full API documentation.

References
----------

- A. Guttman.  R-trees: A Dynamic Index Structure for Spatial Searching.
  Proceedings of ACM SIGMOD, pages 47-57, 1984.
  http://www.cs.jhu.edu/~misha/ReadingSeminar/Papers/Guttman84.pdf

- N. Beckmann, H .P. Kriegel, R. Schneider and B. Seeger.  The R*-tree: An
  Efficient and Robust Access Method for Points and Rectangles.  Proceedings
  of ACM SIGMOD, pages 323-331, May 1990.
  http://infolab.usc.edu/csci587/Fall2011/papers/p322-beckmann.pdf

- N. Roussopoulos, S. Kelley and F. Vincent.  Nearest Neighbor Queries.  ACM
  SIGMOD, pages 71-79, 1995.
  http://www.postgis.org/support/nearestneighbor.pdf

Author
------

Written by [Daniel Connelly](http://dhconnelly.com) (<dhconnelly@gmail.com>).

License
-------

rtreego is released under a BSD-style license, described here and in the
`LICENSE` file:

Copyright (c) 2012, Daniel Connelly. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of Daniel Connelly nor the names of its contributors may be
   used to endorse or promote products derived from this software without
   specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
