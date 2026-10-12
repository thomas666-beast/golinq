// Package linq provides LINQ-style lazy query operators for Go,
// built on iter.Seq.
package linq

import (
	"iter"
	"slices"
)

// Query is a lazy sequence of T. Operators build new queries without
// doing any work until the query is enumerated.
type Query[T any] struct {
	seq iter.Seq[T]
}

// From creates a query over a slice. The slice is read each time the
// query is enumerated, so later changes to it are visible.
func From[T any](source []T) Query[T] {
	return Query[T]{seq: slices.Values(source)}
}

// FromSeq creates a query over any iter.Seq.
// A one-shot sequence (such as one reading a channel) can only be
// enumerated once.
func FromSeq[T any](seq iter.Seq[T]) Query[T] {
	return Query[T]{seq: seq}
}

// Seq returns the underlying sequence, so a query works with
// `for x := range q.Seq()` and the standard library.
func (q Query[T]) Seq() iter.Seq[T] {
	if q.seq == nil {
		return func(yield func(T) bool) {}
	}
	return q.seq
}

// ToSlice runs the query and collects the results into a slice.
// The result is never nil.
func (q Query[T]) ToSlice() []T {
	result := slices.Collect(q.Seq())
	if result == nil {
		return []T{}
	}
	return result
}

// Where keeps only the elements for which predicate returns true.
func (q Query[T]) Where(predicate func(T) bool) Query[T] {
	return Query[T]{seq: func(yield func(T) bool) {
		for x := range q.Seq() {
			if predicate(x) && !yield(x) {
				return
			}
		}
	}}
}

// Select transforms each element using selector.
// It is a function rather than a method because Go methods cannot
// introduce new type parameters.
func Select[T, R any](q Query[T], selector func(T) R) Query[R] {
	return Query[R]{seq: func(yield func(R) bool) {
		for x := range q.Seq() {
			if !yield(selector(x)) {
				return
			}
		}
	}}
}
