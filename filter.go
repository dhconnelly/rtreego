package rtreego

// Filter is an interface for filtering leaves during search. The parameters
// should be treated as read-only. If refuse is true, the currenty entry will
// not be added to the result set. If abort is true, the search is aborted and
// the current result set will be returned.
type Filter interface {
	Filter(results []Spatial, object Spatial) (refuse, abort bool)
}

// ApplyFilters applies the given filters and returns their consensus.
func applyFilters(results []Spatial, object Spatial, filters []Filter) (bool, bool) {
	var refuse, abort bool
	for _, f := range filters {
		ref, abt := f.Filter(results, object)

		if ref {
			refuse = true
		}

		if abt {
			abort = true
		}

		// some filter after the aborting filter might still refuse the leaf,
		// so we can only break early if both are true
		if refuse && abort {
			break
		}
	}

	return refuse, abort
}

// LimitFilter aborts the search after certain amount of results have been
// gathered.
type LimitFilter struct {
	limit int
}

// NewLimitFilter returns a new limit filter.
func NewLimitFilter(limit int) *LimitFilter {
	return &LimitFilter{
		limit: limit,
	}
}

// Filter checks if the results have reached the limit size and aborts if so.
func (f *LimitFilter) Filter(results []Spatial, object Spatial) (bool, bool) {
	if len(results) >= f.limit {
		return true, true
	}

	return false, false
}
