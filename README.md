rtreego
=======

rtreego is a library for efficiently storing and querying spatial data in the
[Go programming language](http://golang.org).

Overview
--------

The [R-tree](http://en.wikipedia.org/wiki/R-tree) is a popular data structure for efficiently storing and querying spatial objects [Gut84].  It is often used for implementing geospatial indexes in database systems.  The variant implemented here, known as the R*-tree, improves performance and increases storage utilization with little implementation overhead, in addition to offering significant improvements when handling point objects [Beck90].  In addition to bounding box queries, k-nearest-neighbor queries are also supported [Rous95].

R*-trees are balanced, so maximum tree height is guaranteed to be logarithmic in the number of entries; however, there is no guarantee of good query performance.  Instead, a number of heuristics are applied that perform well in practice.  For more details please refer to the references.

Status
------

rtreego is currently in the initial stages of development and is not ready for use.

Usage
-----

Installation
------------

Assuming you're using a recent weekly build of Go (at least weekly.2012-01-27 11507), `go install github.com/dhconnelly/rtreego`.  Then `import "github.com/dhconnelly/rtreego"` in your source files.

References
----------

- A. Guttman. [R-trees: A Dynamic Index Structure for Spatial Searching](http://www.cs.jhu.edu/~misha/ReadingSeminar/Papers/Guttman84.pdf). Proceedings of ACM SIGMOD, pages 47-57, 1984.
- N. Beckmann, H .P. Kriegel, R. Schneider and B. Seeger. [The R*-tree: An Efficient and Robust Access Method for Points and Rectangles](http://infolab.usc.edu/csci587/Fall2011/papers/p322-beckmann.pdf). Proceedings of ACM SIGMOD, pages 323-331, May 1990.
- N. Roussopoulos, S. Kelley and F. Vincent. [Nearest Neighbor Queries](http://www.postgis.org/support/nearestneighbor.pdf). ACM SIGMOD, pages 71-79, 1995.

About
-----

rtreego is written and maintained by [Daniel Connelly](http://dhconnelly.com) and is released under a BSD-style license.  See the LICENSE file for more details.
