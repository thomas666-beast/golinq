package linq_test

import (
	"errors"
	"maps"
	"slices"
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

// Zip
func TestZip(t *testing.T) {
	got := linq.Zip(
		linq.From([]int{1, 2}), linq.From([]int{10, 20}),
		func(a, b int) int { return a + b },
	).ToSlice()
	if !slices.Equal(got, []int{11, 22}) {
		t.Fatalf("got %v", got)
	}
}

func TestZipStopsAtShortest(t *testing.T) {
	pair := func(a int, b string) string { return b }
	if got := linq.Zip(linq.From([]int{1, 2, 3}), linq.From([]string{"a"}), pair).ToSlice(); !slices.Equal(got, []string{"a"}) {
		t.Fatalf("got %v", got)
	}
	if got := linq.Zip(linq.From([]int{1}), linq.From([]string{"a", "b"}), pair).ToSlice(); !slices.Equal(got, []string{"a"}) {
		t.Fatalf("got %v", got)
	}
}

func TestZipPairs(t *testing.T) {
	got := linq.ZipPairs(linq.From([]int{1, 2}), linq.From([]string{"a", "b"})).ToSlice()
	want := []linq.Pair[int, string]{{First: 1, Second: "a"}, {First: 2, Second: "b"}}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v", got)
	}
}

func TestZipEmpty(t *testing.T) {
	got := linq.ZipPairs(linq.From([]int{}), linq.From([]int{1, 2})).ToSlice()
	if len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestZipWorksWithInfiniteSequence(t *testing.T) {
	got := linq.Zip(
		infiniteFrom1(), linq.From([]string{"a", "b"}),
		func(n int, s string) string { return s },
	).ToSlice()
	if !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("got %v", got)
	}
}

func TestZipStopsPullingOnEarlyBreak(t *testing.T) {
	var pulled []int
	q := linq.ZipPairs(trackedSource(&pulled), linq.From([]int{1, 2, 3, 4}))
	for range q.Seq() {
		break
	}
	if !slices.Equal(pulled, []int{0}) {
		t.Fatalf("pulled %v", pulled)
	}
}

// Concat
func TestConcat(t *testing.T) {
	got := linq.From([]int{1, 2}).Concat(linq.From([]int{3, 4})).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Fatalf("got %v", got)
	}
}

func TestConcatKeepsDuplicates(t *testing.T) {
	got := linq.From([]int{1, 2}).Concat(linq.From([]int{2, 3})).ToSlice()
	if !slices.Equal(got, []int{1, 2, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestConcatWithInfiniteFirst(t *testing.T) {
	got := infiniteFrom1().Concat(linq.From([]int{99})).Take(3).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestConcatIsReEnumerable(t *testing.T) {
	q := linq.From([]int{1}).Concat(linq.From([]int{2}))
	if !slices.Equal(q.ToSlice(), []int{1, 2}) || !slices.Equal(q.ToSlice(), []int{1, 2}) {
		t.Fatal("enumerations differ")
	}
}

// Reverse
func TestReverse(t *testing.T) {
	got := linq.From([]int{1, 2, 3}).Reverse().ToSlice()
	if !slices.Equal(got, []int{3, 2, 1}) {
		t.Fatalf("got %v", got)
	}
}

func TestReverseEmpty(t *testing.T) {
	if got := linq.From([]int{}).Reverse().ToSlice(); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestReverseAfterOrderBy(t *testing.T) {
	got := linq.OrderBy(linq.From([]int{2, 3, 1}), ident).Reverse().ToSlice()
	if !slices.Equal(got, []int{3, 2, 1}) {
		t.Fatalf("got %v", got)
	}
}

func TestReverseDeferredExecution(t *testing.T) {
	var calls []int
	q := linq.Select(linq.From([]int{1, 2, 3}), func(x int) int {
		calls = append(calls, x)
		return x
	}).Reverse()
	if len(calls) != 0 {
		t.Fatalf("ran before enumeration: %v", calls)
	}
	if got := q.ToSlice(); !slices.Equal(got, []int{3, 2, 1}) {
		t.Fatalf("got %v", got)
	}
	if !slices.Equal(calls, []int{1, 2, 3}) {
		t.Fatalf("calls %v", calls)
	}
}

func TestReverseDoesNotMutateSource(t *testing.T) {
	data := []int{1, 2, 3}
	linq.From(data).Reverse().ToSlice()
	if !slices.Equal(data, []int{1, 2, 3}) {
		t.Fatalf("source was modified: %v", data)
	}
}

// ToMap
func TestToMap(t *testing.T) {
	got, err := linq.ToMap(linq.From([]string{"a", "bb", "ccc"}), func(s string) int { return len(s) })
	want := map[int]string{1: "a", 2: "bb", 3: "ccc"}
	if err != nil || !maps.Equal(got, want) {
		t.Fatalf("got %v, %v", got, err)
	}
}

func TestToMapElement(t *testing.T) {
	got, err := linq.ToMapElement(
		linq.From([]string{"a", "bb"}),
		func(s string) string { return s },
		func(s string) int { return len(s) },
	)
	if err != nil || !maps.Equal(got, map[string]int{"a": 1, "bb": 2}) {
		t.Fatalf("got %v, %v", got, err)
	}
}

func TestToMapDuplicateKey(t *testing.T) {
	got, err := linq.ToMap(linq.From([]string{"a", "b"}), func(s string) int { return len(s) })
	if !errors.Is(err, linq.ErrDuplicateKey) || got != nil {
		t.Fatalf("got %v, %v", got, err)
	}
}

func TestToMapEmpty(t *testing.T) {
	got, err := linq.ToMap(linq.From([]int{}), ident)
	if err != nil || len(got) != 0 {
		t.Fatalf("got %v, %v", got, err)
	}
}

// ToSet
func TestToSet(t *testing.T) {
	got := linq.ToSet(linq.From([]int{1, 2, 2, 3}))
	if !maps.Equal(got, map[int]struct{}{1: {}, 2: {}, 3: {}}) {
		t.Fatalf("got %v", got)
	}
}

func TestToSetEmpty(t *testing.T) {
	if got := linq.ToSet(linq.From([]int{})); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}
