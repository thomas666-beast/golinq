package linq

import "cmp"

// Number is satisfied by all integer and floating-point types,
// including named types derived from them.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Count returns the number of elements.
func (q Query[T]) Count() int {
	n := 0
	for range q.Seq() {
		n++
	}
	return n
}

// CountWhere returns the number of elements satisfying predicate.
func (q Query[T]) CountWhere(predicate func(T) bool) int {
	n := 0
	for x := range q.Seq() {
		if predicate(x) {
			n++
		}
	}
	return n
}

// Sum adds up the elements. An empty query gives 0.
func Sum[N Number](q Query[N]) N {
	var total N
	for x := range q.Seq() {
		total += x
	}
	return total
}

// SumBy adds up the values returned by selector.
func SumBy[T any, N Number](q Query[T], selector func(T) N) N {
	var total N
	for x := range q.Seq() {
		total += selector(x)
	}
	return total
}

// Min returns the smallest element. ok is false if the query is empty.
func Min[T cmp.Ordered](q Query[T]) (value T, ok bool) {
	for x := range q.Seq() {
		if !ok || x < value {
			value, ok = x, true
		}
	}
	return value, ok
}

// Max returns the largest element. ok is false if the query is empty.
func Max[T cmp.Ordered](q Query[T]) (value T, ok bool) {
	for x := range q.Seq() {
		if !ok || x > value {
			value, ok = x, true
		}
	}
	return value, ok
}

// MinBy returns the element with the smallest key. If several elements
// share the smallest key, the first one wins.
func MinBy[T any, K cmp.Ordered](q Query[T], key func(T) K) (value T, ok bool) {
	var best K
	for x := range q.Seq() {
		k := key(x)
		if !ok || k < best {
			value, best, ok = x, k, true
		}
	}
	return value, ok
}

// MaxBy returns the element with the largest key. If several elements
// share the largest key, the first one wins.
func MaxBy[T any, K cmp.Ordered](q Query[T], key func(T) K) (value T, ok bool) {
	var best K
	for x := range q.Seq() {
		k := key(x)
		if !ok || k > best {
			value, best, ok = x, k, true
		}
	}
	return value, ok
}

// Average returns the mean of the elements. ok is false if the query
// is empty.
func Average[N Number](q Query[N]) (float64, bool) {
	return AverageBy(q, func(x N) N { return x })
}

// AverageBy returns the mean of the values returned by selector.
func AverageBy[T any, N Number](q Query[T], selector func(T) N) (float64, bool) {
	var total float64
	n := 0
	for x := range q.Seq() {
		total += float64(selector(x))
		n++
	}
	if n == 0 {
		return 0, false
	}
	return total / float64(n), true
}

// Aggregate folds the query into one value, starting from seed.
// An empty query returns seed. The accumulator type may differ from
// the element type.
func Aggregate[T, A any](q Query[T], seed A, fn func(acc A, x T) A) A {
	acc := seed
	for x := range q.Seq() {
		acc = fn(acc, x)
	}
	return acc
}

// Reduce folds the query using the first element as the starting
// accumulator. ok is false if the query is empty.
func (q Query[T]) Reduce(fn func(acc, x T) T) (value T, ok bool) {
	for x := range q.Seq() {
		if !ok {
			value, ok = x, true
			continue
		}
		value = fn(value, x)
	}
	return value, ok
}
