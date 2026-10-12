package linq

import "errors"

var (
	// ErrNoElements is returned by Single when no element matches.
	ErrNoElements = errors.New("linq: sequence contains no matching element")
	// ErrMultipleElements is returned by Single when more than one element matches.
	ErrMultipleElements = errors.New("linq: sequence contains more than one matching element")
)

// First returns the first element. ok is false if the query is empty.
func (q Query[T]) First() (value T, ok bool) {
	for x := range q.Seq() {
		return x, true
	}
	return value, false
}

// FirstWhere returns the first element matching predicate.
func (q Query[T]) FirstWhere(predicate func(T) bool) (value T, ok bool) {
	for x := range q.Seq() {
		if predicate(x) {
			return x, true
		}
	}
	return value, false
}

// Single returns the only element. It returns ErrNoElements if the query
// is empty and ErrMultipleElements if it has more than one element.
func (q Query[T]) Single() (T, error) {
	return q.SingleWhere(func(T) bool { return true })
}

// SingleWhere returns the only element matching predicate. It stops as
// soon as it finds a second match.
func (q Query[T]) SingleWhere(predicate func(T) bool) (T, error) {
	var result T
	found := false
	for x := range q.Seq() {
		if !predicate(x) {
			continue
		}
		if found {
			var zero T
			return zero, ErrMultipleElements
		}
		found = true
		result = x
	}
	if !found {
		return result, ErrNoElements
	}
	return result, nil
}

// Last returns the last element. ok is false if the query is empty.
func (q Query[T]) Last() (value T, ok bool) {
	for x := range q.Seq() {
		value, ok = x, true
	}
	return value, ok
}

// LastWhere returns the last element matching predicate.
func (q Query[T]) LastWhere(predicate func(T) bool) (value T, ok bool) {
	for x := range q.Seq() {
		if predicate(x) {
			value, ok = x, true
		}
	}
	return value, ok
}

// ElementAt returns the element at a zero-based index. ok is false if
// the index is negative or past the end.
func (q Query[T]) ElementAt(index int) (value T, ok bool) {
	if index < 0 {
		return value, false
	}
	i := 0
	for x := range q.Seq() {
		if i == index {
			return x, true
		}
		i++
	}
	return value, false
}

// Any reports whether the query has at least one element.
func (q Query[T]) Any() bool {
	for range q.Seq() {
		return true
	}
	return false
}

// AnyWhere reports whether any element satisfies predicate.
func (q Query[T]) AnyWhere(predicate func(T) bool) bool {
	for x := range q.Seq() {
		if predicate(x) {
			return true
		}
	}
	return false
}

// All reports whether every element satisfies predicate.
// It is true for an empty query.
func (q Query[T]) All(predicate func(T) bool) bool {
	for x := range q.Seq() {
		if !predicate(x) {
			return false
		}
	}
	return true
}

// Contains reports whether the query contains value.
// It is a function because it needs T to be comparable.
func Contains[T comparable](q Query[T], value T) bool {
	for x := range q.Seq() {
		if x == value {
			return true
		}
	}
	return false
}
