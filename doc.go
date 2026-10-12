// Package linq provides LINQ-style lazy query operators for Go,
// built on iter.Seq.
//
// Wrap a slice or sequence with From or FromSeq, chain operators, and
// finish with a terminal operation such as ToSlice, First, Count or Sum:
//
//	adults := linq.From(people).
//		Where(func(p Person) bool { return p.Age >= 18 }).
//		ToSlice()
//
// # Methods and functions
//
// Go methods cannot declare their own type parameters, so operators that
// keep the element type (Where, Take, Skip, Concat, Reverse, ...) are
// methods, and operators that change the type or need a constraint
// (Select, GroupBy, Join, Distinct, Sum, OrderBy, ...) are functions:
//
//	names := linq.Select(
//		linq.From(people).Where(isAdult),
//		func(p Person) string { return p.Name },
//	).ToSlice()
//
// # Deferred execution
//
// Operators that return a Query do no work until it is enumerated, and a
// query can be enumerated more than once. Each enumeration re-runs it.
// The exception is a query built with FromSeq on a one-shot sequence,
// which can only be consumed once. Operators that return a value
// (First, Count, Sum, ToSlice, ...) run immediately.
//
// # Errors
//
// Operators do not panic on empty input. Where C# would throw, linq
// returns an ok bool or an error (see ErrNoElements, ErrMultipleElements
// and ErrDuplicateKey). The one exception is Chunk, which panics when
// size is less than 1, like slices.Chunk.
package linq
