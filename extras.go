package linq

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

// ErrDuplicateKey is returned by ToMap and ToMapElement when two
// elements produce the same key.
var ErrDuplicateKey = errors.New("linq: duplicate key")

// Pair holds one element from each of two queries. See ZipPairs.
type Pair[A, B any] struct {
	First  A
	Second B
}

// Concat appends second to the query. Duplicates are kept.
func (q Query[T]) Concat(second Query[T]) Query[T] {
	return Query[T]{seq: func(yield func(T) bool) {
		for x := range q.Seq() {
			if !yield(x) {
				return
			}
		}
		for x := range second.Seq() {
			if !yield(x) {
				return
			}
		}
	}}
}

// Reverse reverses the order of the elements. The whole query is
// buffered when the result is enumerated.
func (q Query[T]) Reverse() Query[T] {
	return Query[T]{seq: func(yield func(T) bool) {
		items := slices.Collect(q.Seq())
		for i := len(items) - 1; i >= 0; i-- {
			if !yield(items[i]) {
				return
			}
		}
	}}
}

// Zip pairs elements from both queries using result, stopping at the
// shorter one.
func Zip[T, U, R any](first Query[T], second Query[U], result func(T, U) R) Query[R] {
	return Query[R]{seq: func(yield func(R) bool) {
		next, stop := iter.Pull(second.Seq())
		defer stop()
		for a := range first.Seq() {
			b, ok := next()
			if !ok {
				return
			}
			if !yield(result(a, b)) {
				return
			}
		}
	}}
}

// ZipPairs is like Zip but yields Pair values.
func ZipPairs[T, U any](first Query[T], second Query[U]) Query[Pair[T, U]] {
	return Zip(first, second, func(a T, b U) Pair[T, U] {
		return Pair[T, U]{First: a, Second: b}
	})
}

// ToMap builds a map keyed by key. It returns ErrDuplicateKey (wrapped
// with the key) if two elements produce the same key.
func ToMap[T any, K comparable](q Query[T], key func(T) K) (map[K]T, error) {
	return ToMapElement(q, key, func(x T) T { return x })
}

// ToMapElement is like ToMap but stores the value returned by element.
func ToMapElement[T any, K comparable, V any](
	q Query[T],
	key func(T) K,
	element func(T) V,
) (map[K]V, error) {
	result := make(map[K]V)
	for x := range q.Seq() {
		k := key(x)
		if _, exists := result[k]; exists {
			return nil, fmt.Errorf("%w: %v", ErrDuplicateKey, k)
		}
		result[k] = element(x)
	}
	return result, nil
}

// ToSet collects the elements into a set. It is a function because it
// needs T to be comparable.
func ToSet[T comparable](q Query[T]) map[T]struct{} {
	result := make(map[T]struct{})
	for x := range q.Seq() {
		result[x] = struct{}{}
	}
	return result
}
