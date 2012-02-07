rtreego
=======

rtreego is a library for efficiently storing and querying spatial data in the
[Go programming language](http://golang.org).

Overview
--------

The classic R-tree is a popular and efficient data structure for storing and querying spatial objects [Gut84].  The variant implemented here, known as the R*-tree, improves performance and increases storage utilization with little implementation overhead, in addition to offering significant improvements when handling point objects [Beck90].  In addition to rectangle intersection queries (ie, bounding box queries), k-nearest-neighbor queries are also supported [Rous95].

The R*-tree is a dynamic data structure--insertions and deletions can be performed in any order.  Further, insertions and deletions trigger rebalancing automatically, so that no explicit rebalancing is required by the user.  Maximum tree height is guaranteed to be logarithmic in the number of entries, but there is no guarantee of good worst-case query performance; a number of heuristics are applied that perform well in practice.  For more details on these heuristics please refer to the references.

Usage
-----


Installation
------------

Using a recently updated Go programming language (at least weekly.2012-01-27 11507), simply

`go install github.com/dhconnelly/rtreego`.

Then `import "github.com/dhconnelly/rtreego"` in your source files.

References
----------

- A. Guttman, [R-trees: A Dynamic Index Structure for Spatial Searching](http://www.cs.jhu.edu/~misha/ReadingSeminar/Papers/Guttman84.pdf). Proceedings of ACM SIGMOD, pages 47-57, 1984.
- N. Beckmann, H .P. Kriegel, R. Schneider and B. Seeger. [The R*-tree: An Efficient and Robust Access Method for Points and Rectangles](http://infolab.usc.edu/csci587/Fall2011/papers/p322-beckmann.pdf). Proceedings of ACM SIGMOD, pages 323-331, May 1990.
- N. Roussopoulos, S. Kelley and F. Vincent, [Nearest Neighbor Queries](http://www.postgis.org/support/nearestneighbor.pdf). ACM SIGMOD, pages 71-79, 1995.

About
-----

This library was written by [Daniel Connelly](http://dhconnelly.com) and is released under a BSD-style license.  See the LICENSE file for more details.
