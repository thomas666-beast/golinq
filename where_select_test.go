package linq_test

import (
	"slices"
	"strconv"
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

func TestWhereFilters(t *testing.T) {
	got := linq.From([]int{1, 2, 3, 4, 5}).
		Where(func(x int) bool { return x%2 == 0 }).
		ToSlice()
	if !slices.Equal(got, []int{2, 4}) {
		t.Fatalf("got %v", got)
	}
}

func TestWhereNoMatches(t *testing.T) {
	got := linq.From([]int{1, 2, 3}).
		Where(func(x int) bool { return x > 10 }).
		ToSlice()
	if len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestSelectTransforms(t *testing.T) {
	got := linq.Select(linq.From([]int{1, 2, 3}), func(x int) int { return x * 10 }).ToSlice()
	if !slices.Equal(got, []int{10, 20, 30}) {
		t.Fatalf("got %v", got)
	}
}

func TestSelectChangesType(t *testing.T) {
	got := linq.Select(linq.From([]int{1, 2, 3}), strconv.Itoa).ToSlice()
	if !slices.Equal(got, []string{"1", "2", "3"}) {
		t.Fatalf("got %v", got)
	}
}

func TestChaining(t *testing.T) {
	evens := linq.From([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).
		Where(func(x int) bool { return x%2 == 0 })
	got := linq.Select(evens, func(x int) int { return x * x }).ToSlice()
	if !slices.Equal(got, []int{0, 4, 16, 36, 64}) {
		t.Fatalf("got %v", got)
	}
}

func TestDeferredExecution(t *testing.T) {
	var calls []int
	q := linq.Select(linq.From([]int{1, 2, 3}), func(x int) int {
		calls = append(calls, x)
		return x
	})
	if len(calls) != 0 {
		t.Fatalf("selector ran before enumeration: %v", calls)
	}
	q.ToSlice()
	if !slices.Equal(calls, []int{1, 2, 3}) {
		t.Fatalf("got %v", calls)
	}
}

func TestReEnumerationRerunsSelector(t *testing.T) {
	var calls []int
	q := linq.Select(linq.From([]int{1, 2}), func(x int) int {
		calls = append(calls, x)
		return x
	})
	q.ToSlice()
	q.ToSlice()
	if !slices.Equal(calls, []int{1, 2, 1, 2}) {
		t.Fatalf("got %v", calls)
	}
}

func TestWhereSelectStopOnEarlyBreak(t *testing.T) {
	var pulled []int
	source := func(yield func(int) bool) {
		for i := 0; i < 100; i++ {
			pulled = append(pulled, i)
			if !yield(i) {
				return
			}
		}
	}
	q := linq.Select(
		linq.FromSeq(source).Where(func(x int) bool { return true }),
		func(x int) int { return x },
	)
	for x := range q.Seq() {
		if x == 2 {
			break
		}
	}
	if !slices.Equal(pulled, []int{0, 1, 2}) {
		t.Fatalf("pulled %v", pulled)
	}
}
