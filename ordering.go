package linq

import (
	"cmp"
	"slices"
)

// OrderedQuery is a Query with a sort order. Use ThenBy and
// ThenByDescending to add secondary keys. Because it embeds Query,
// all Query methods are available on it. To pass it to a function
// such as Select, use its Query field.
type OrderedQuery[T any] struct {
	Query[T]
	source    Query[T]
	comparers []func(a, b T) int
}

func keyComparer[T any, K cmp.Ordered](key func(T) K, descending bool) func(a, b T) int {
	if descending {
		return func(a, b T) int { return cmp.Compare(key(b), key(a)) }
	}
	return func(a, b T) int { return cmp.Compare(key(a), key(b)) }
}

func newOrdered[T any](source Query[T], comparers []func(a, b T) int) OrderedQuery[T] {
	seq := func(yield func(T) bool) {
		items := slices.Collect(source.Seq()) // a copy, so the source is untouched
		slices.SortStableFunc(items, func(a, b T) int {
			for _, c := range comparers {
				if r := c(a, b); r != 0 {
					return r
				}
			}
			return 0
		})
		for _, x := range items {
			if !yield(x) {
				return
			}
		}
	}
	return OrderedQuery[T]{
		Query:     Query[T]{seq: seq},
		source:    source,
		comparers: comparers,
	}
}

// OrderBy sorts ascending by the key returned from key.
// The sort is stable. It runs when the query is enumerated.
func OrderBy[T any, K cmp.Ordered](q Query[T], key func(T) K) OrderedQuery[T] {
	return newOrdered(q, []func(a, b T) int{keyComparer(key, false)})
}

// OrderByDescending sorts descending by the key returned from key.
func OrderByDescending[T any, K cmp.Ordered](q Query[T], key func(T) K) OrderedQuery[T] {
	return newOrdered(q, []func(a, b T) int{keyComparer(key, true)})
}

// ThenBy adds a secondary ascending sort key.
func ThenBy[T any, K cmp.Ordered](oq OrderedQuery[T], key func(T) K) OrderedQuery[T] {
	comparers := append(slices.Clone(oq.comparers), keyComparer(key, false))
	return newOrdered(oq.source, comparers)
}

// ThenByDescending adds a secondary descending sort key.
func ThenByDescending[T any, K cmp.Ordered](oq OrderedQuery[T], key func(T) K) OrderedQuery[T] {
	comparers := append(slices.Clone(oq.comparers), keyComparer(key, true))
	return newOrdered(oq.source, comparers)
}
