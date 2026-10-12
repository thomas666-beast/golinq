package linq

// Grouping is a group of elements sharing a common key. Because it
// embeds Query, all Query methods work on it. To pass it to a function
// such as Select or Sum, use its Query field.
type Grouping[K comparable, T any] struct {
	Query[T]
	Key K
}

// GroupBy groups elements by key. Groups are yielded in order of first
// appearance, and elements keep their original order within a group.
// The whole source is read when the first group is requested.
func GroupBy[T any, K comparable](q Query[T], key func(T) K) Query[Grouping[K, T]] {
	return GroupByElement(q, key, func(x T) T { return x })
}

// GroupByElement groups elements by key, storing the value returned
// by element in each group instead of the original element.
func GroupByElement[T any, K comparable, E any](
	q Query[T],
	key func(T) K,
	element func(T) E,
) Query[Grouping[K, E]] {
	return Query[Grouping[K, E]]{seq: func(yield func(Grouping[K, E]) bool) {
		var order []K
		groups := make(map[K][]E)
		for x := range q.Seq() {
			k := key(x)
			if _, exists := groups[k]; !exists {
				order = append(order, k)
			}
			groups[k] = append(groups[k], element(x))
		}
		for _, k := range order {
			g := Grouping[K, E]{Query: From(groups[k]), Key: k}
			if !yield(g) {
				return
			}
		}
	}}
}

func buildLookup[U any, K comparable](inner Query[U], key func(U) K) map[K][]U {
	lookup := make(map[K][]U)
	for x := range inner.Seq() {
		k := key(x)
		lookup[k] = append(lookup[k], x)
	}
	return lookup
}

// Join is an inner join: it pairs each outer element with every inner
// element that has the same key. Outer elements with no match are
// dropped. inner is read completely when the result is enumerated.
func Join[T, U any, K comparable, R any](
	outer Query[T],
	inner Query[U],
	outerKey func(T) K,
	innerKey func(U) K,
	result func(T, U) R,
) Query[R] {
	return Query[R]{seq: func(yield func(R) bool) {
		lookup := buildLookup(inner, innerKey)
		for o := range outer.Seq() {
			for _, i := range lookup[outerKey(o)] {
				if !yield(result(o, i)) {
					return
				}
			}
		}
	}}
}

// GroupJoin is a left join: every outer element appears once, paired
// with a Query of its matching inner elements (empty if none).
func GroupJoin[T, U any, K comparable, R any](
	outer Query[T],
	inner Query[U],
	outerKey func(T) K,
	innerKey func(U) K,
	result func(T, Query[U]) R,
) Query[R] {
	return Query[R]{seq: func(yield func(R) bool) {
		lookup := buildLookup(inner, innerKey)
		for o := range outer.Seq() {
			if !yield(result(o, From(lookup[outerKey(o)]))) {
				return
			}
		}
	}}
}
