package linq_test

import (
	"slices"
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

// Distinct / DistinctBy
func TestDistinct(t *testing.T) {
	got := linq.Distinct(linq.From([]int{1, 2, 2, 3, 1})).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestDistinctPreservesFirstOccurrenceOrder(t *testing.T) {
	got := linq.Distinct(linq.From([]int{3, 1, 3, 2, 1})).ToSlice()
	if !slices.Equal(got, []int{3, 1, 2}) {
		t.Fatalf("got %v", got)
	}
}

func TestDistinctEmpty(t *testing.T) {
	if got := linq.Distinct(linq.From([]int{})).ToSlice(); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestDistinctBy(t *testing.T) {
	words := []string{"apple", "avocado", "banana", "blueberry", "cherry"}
	got := linq.DistinctBy(linq.From(words), func(w string) byte { return w[0] }).ToSlice()
	if !slices.Equal(got, []string{"apple", "banana", "cherry"}) {
		t.Fatalf("got %v", got)
	}
}

func TestDistinctByDeferredExecution(t *testing.T) {
	var calls []int
	q := linq.DistinctBy(linq.From([]int{1, 1, 2}), func(x int) int {
		calls = append(calls, x)
		return x
	})
	if len(calls) != 0 {
		t.Fatalf("key ran before enumeration: %v", calls)
	}
	if got := q.ToSlice(); !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("got %v", got)
	}
	if !slices.Equal(calls, []int{1, 1, 2}) {
		t.Fatalf("calls %v", calls)
	}
}

func TestDistinctWorksOnInfiniteSequence(t *testing.T) {
	got := linq.Distinct(infiniteFrom1()).Take(3).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

// Union
func TestUnion(t *testing.T) {
	got := linq.Union(linq.From([]int{1, 2, 3}), linq.From([]int{3, 4, 5})).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("got %v", got)
	}
}

func TestUnionRemovesDuplicatesWithinEachSide(t *testing.T) {
	got := linq.Union(linq.From([]int{1, 1, 2}), linq.From([]int{2, 3, 3})).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestUnionOrder(t *testing.T) {
	got := linq.Union(linq.From([]int{5, 1}), linq.From([]int{3, 1})).ToSlice()
	if !slices.Equal(got, []int{5, 1, 3}) {
		t.Fatalf("got %v", got)
	}
}

// Intersect
func TestIntersect(t *testing.T) {
	got := linq.Intersect(linq.From([]int{1, 2, 3, 4}), linq.From([]int{3, 4, 5})).ToSlice()
	if !slices.Equal(got, []int{3, 4}) {
		t.Fatalf("got %v", got)
	}
}

func TestIntersectRemovesDuplicates(t *testing.T) {
	got := linq.Intersect(linq.From([]int{1, 1, 2, 2}), linq.From([]int{1, 2})).ToSlice()
	if !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("got %v", got)
	}
}

func TestIntersectFollowsFirstQueryOrder(t *testing.T) {
	got := linq.Intersect(linq.From([]int{3, 1, 2}), linq.From([]int{1, 2, 3})).ToSlice()
	if !slices.Equal(got, []int{3, 1, 2}) {
		t.Fatalf("got %v", got)
	}
}

func TestIntersectNoOverlap(t *testing.T) {
	got := linq.Intersect(linq.From([]int{1, 2}), linq.From([]int{3, 4})).ToSlice()
	if len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestIntersectWorksOnInfiniteFirst(t *testing.T) {
	got := linq.Intersect(infiniteFrom1(), linq.From([]int{2, 4})).Take(2).ToSlice()
	if !slices.Equal(got, []int{2, 4}) {
		t.Fatalf("got %v", got)
	}
}

// Except
func TestExcept(t *testing.T) {
	got := linq.Except(linq.From([]int{1, 2, 3, 4}), linq.From([]int{2, 4})).ToSlice()
	if !slices.Equal(got, []int{1, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestExceptRemovesDuplicates(t *testing.T) {
	got := linq.Except(linq.From([]int{1, 1, 2, 3, 3}), linq.From([]int{2})).ToSlice()
	if !slices.Equal(got, []int{1, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestExceptEmptySecond(t *testing.T) {
	got := linq.Except(linq.From([]int{1, 2, 2}), linq.From([]int{})).ToSlice()
	if !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("got %v", got)
	}
}

// General
func TestSetOperationsAreReEnumerable(t *testing.T) {
	a, b := linq.From([]int{1, 2, 3}), linq.From([]int{2, 3, 4})
	queries := []linq.Query[int]{
		linq.Distinct(a),
		linq.Union(a, b),
		linq.Intersect(a, b),
		linq.Except(a, b),
	}
	for i, q := range queries {
		if !slices.Equal(q.ToSlice(), q.ToSlice()) {
			t.Fatalf("query %d differs between enumerations", i)
		}
	}
}

func TestSecondQueryIsReadLazily(t *testing.T) {
	read := false
	second := linq.FromSeq(func(yield func(int) bool) {
		read = true
		yield(2)
	})
	q := linq.Intersect(linq.From([]int{1, 2}), second)
	if read {
		t.Fatal("second query was read before enumeration")
	}
	if got := q.ToSlice(); !slices.Equal(got, []int{2}) {
		t.Fatalf("got %v", got)
	}
	if !read {
		t.Fatal("second query was never read")
	}
}
