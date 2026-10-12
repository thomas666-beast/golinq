package linq

// Distinct removes duplicates, keeping the first occurrence of each.
// It is a function because it needs T to be comparable.
func Distinct[T comparable](q Query[T]) Query[T] {
	return DistinctBy(q, func(x T) T { return x })
}

// DistinctBy removes elements whose key was already seen, keeping the
// first occurrence of each key.
func DistinctBy[T any, K comparable](q Query[T], key func(T) K) Query[T] {
	return Query[T]{seq: func(yield func(T) bool) {
		seen := make(map[K]struct{})
		for x := range q.Seq() {
			k := key(x)
			if _, dup := seen[k]; dup {
				continue
			}
			seen[k] = struct{}{}
			if !yield(x) {
				return
			}
		}
	}}
}

// Union returns the distinct elements of first followed by the
// distinct elements of second that were not already seen.
func Union[T comparable](first, second Query[T]) Query[T] {
	return Query[T]{seq: func(yield func(T) bool) {
		seen := make(map[T]struct{})
		for _, src := range [2]Query[T]{first, second} {
			for x := range src.Seq() {
				if _, dup := seen[x]; dup {
					continue
				}
				seen[x] = struct{}{}
				if !yield(x) {
					return
				}
			}
		}
	}}
}

// Intersect returns the distinct elements that appear in both queries,
// in the order of first. second is read completely when the result is
// enumerated.
func Intersect[T comparable](first, second Query[T]) Query[T] {
	return Query[T]{seq: func(yield func(T) bool) {
		allowed := make(map[T]struct{})
		for x := range second.Seq() {
			allowed[x] = struct{}{}
		}
		for x := range first.Seq() {
			if _, ok := allowed[x]; !ok {
				continue
			}
			delete(allowed, x) // each value is yielded only once
			if !yield(x) {
				return
			}
		}
	}}
}

// Except returns the distinct elements of first that do not appear in
// second. second is read completely when the result is enumerated.
func Except[T comparable](first, second Query[T]) Query[T] {
	return Query[T]{seq: func(yield func(T) bool) {
		seen := make(map[T]struct{})
		for x := range second.Seq() {
			seen[x] = struct{}{}
		}
		for x := range first.Seq() {
			if _, dup := seen[x]; dup {
				continue
			}
			seen[x] = struct{}{}
			if !yield(x) {
				return
			}
		}
	}}
}
