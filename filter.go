package rtreego

// Filter is an interface for filtering leaves during search. The parameters
// should be treated as read-only. If refuse is true, the currenty entry will
// not be added to the result set. If abort is true, the search is aborted and
// the current result set will be returned.
type Filter func(results []Spatial, object Spatial) (refuse, abort bool)

// ApplyFilters applies the given filters and returns their consensus.
func applyFilters(results []Spatial, object Spatial, filters []Filter) (bool, bool) {
	var refuse, abort bool
	for _, filter := range filters {
		ref, abt := filter(results, object)

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

// LimitFilter checks if the results have reached the limit size and aborts if so.
func LimitFilter(limit int) Filter {
	return func(results []Spatial, object Spatial) (refuse, abort bool) {
		if len(results) >= limit {
			return true, true
		}

		return false, false
	}
}
