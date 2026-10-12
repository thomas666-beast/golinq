package linq_test

import (
	"slices"
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

type person struct {
	Name string
	Pets []string
}

func TestSelectManyFlattens(t *testing.T) {
	got := linq.SelectMany(
		linq.From([][]int{{1, 2}, {3}, {4, 5}}),
		func(x []int) []int { return x },
	).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("got %v", got)
	}
}

func TestSelectManyExtractsInnerCollection(t *testing.T) {
	people := []person{
		{"Ann", []string{"cat", "dog"}},
		{"Bob", nil},
		{"Cy", []string{"fish"}},
	}
	got := linq.SelectMany(linq.From(people), func(p person) []string { return p.Pets }).ToSlice()
	if !slices.Equal(got, []string{"cat", "dog", "fish"}) {
		t.Fatalf("got %v", got)
	}
}

func TestSelectManyResult(t *testing.T) {
	people := []person{
		{"Ann", []string{"cat", "dog"}},
		{"Cy", []string{"fish"}},
	}
	got := linq.SelectManyResult(
		linq.From(people),
		func(p person) []string { return p.Pets },
		func(p person, pet string) string { return p.Name + ":" + pet },
	).ToSlice()
	if !slices.Equal(got, []string{"Ann:cat", "Ann:dog", "Cy:fish"}) {
		t.Fatalf("got %v", got)
	}
}

func TestSelectManyEmptyInner(t *testing.T) {
	got := linq.SelectMany(
		linq.From([][]int{{}, nil}),
		func(x []int) []int { return x },
	).ToSlice()
	if len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestSelectManySeqStreamsInner(t *testing.T) {
	got := linq.SelectManySeq(
		linq.From([]int{1, 2, 3}),
		func(n int) linq.Query[int] {
			return linq.From([]int{n, n * 10}).Where(func(x int) bool { return x != 20 })
		},
	).ToSlice()
	if !slices.Equal(got, []int{1, 10, 2, 3, 30}) {
		t.Fatalf("got %v", got)
	}
}

func TestSelectManyChainsWithWhere(t *testing.T) {
	flat := linq.SelectMany(
		linq.From([][]int{{1, 2, 3}, {4, 5, 6}}),
		func(x []int) []int { return x },
	)
	got := flat.Where(func(n int) bool { return n%2 == 0 }).ToSlice()
	if !slices.Equal(got, []int{2, 4, 6}) {
		t.Fatalf("got %v", got)
	}
}

func TestSelectManyDeferredExecution(t *testing.T) {
	var calls []int
	q := linq.SelectMany(linq.From([]int{1, 2}), func(x int) []int {
		calls = append(calls, x)
		return []int{x, x}
	})
	if len(calls) != 0 {
		t.Fatalf("selector ran before enumeration: %v", calls)
	}
	if got := q.ToSlice(); !slices.Equal(got, []int{1, 1, 2, 2}) {
		t.Fatalf("got %v", got)
	}
	if !slices.Equal(calls, []int{1, 2}) {
		t.Fatalf("got %v", calls)
	}
}

func TestSelectManyStopsOnEarlyBreak(t *testing.T) {
	var calls []int
	q := linq.SelectMany(linq.From([]int{1, 2, 3}), func(x int) []int {
		calls = append(calls, x)
		return []int{x, x}
	})
	count := 0
	for range q.Seq() {
		count++
		if count == 2 { // stop in the middle of the first inner slice batch's end
			break
		}
	}
	if !slices.Equal(calls, []int{1}) {
		t.Fatalf("selector called for %v, expected only [1]", calls)
	}
}
