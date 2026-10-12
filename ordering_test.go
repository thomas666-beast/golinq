package linq_test

import (
	"slices"
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

type emp struct {
	Name string
	Age  int
}

var emps = []emp{{"Bob", 30}, {"Ann", 25}, {"Cy", 30}, {"Dee", 25}}

func ident(x int) int { return x }

func TestOrderByAscending(t *testing.T) {
	got := linq.OrderBy(linq.From([]int{3, 1, 2}), ident).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestOrderByDescending(t *testing.T) {
	got := linq.OrderByDescending(linq.From([]int{3, 1, 2}), ident).ToSlice()
	if !slices.Equal(got, []int{3, 2, 1}) {
		t.Fatalf("got %v", got)
	}
}

func TestOrderByKeySelector(t *testing.T) {
	got := linq.OrderBy(linq.From([]string{"ccc", "a", "bb"}), func(s string) int { return len(s) }).ToSlice()
	if !slices.Equal(got, []string{"a", "bb", "ccc"}) {
		t.Fatalf("got %v", got)
	}
}

func TestOrderByIsStable(t *testing.T) {
	got := linq.OrderBy(linq.From([]string{"bb", "aa", "c", "dd"}), func(s string) int { return len(s) }).ToSlice()
	if !slices.Equal(got, []string{"c", "bb", "aa", "dd"}) {
		t.Fatalf("got %v", got)
	}
}

func TestThenBy(t *testing.T) {
	byAge := linq.OrderBy(linq.From(emps), func(e emp) int { return e.Age })
	got := linq.ThenBy(byAge, func(e emp) string { return e.Name }).ToSlice()
	want := []emp{{"Ann", 25}, {"Dee", 25}, {"Bob", 30}, {"Cy", 30}}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v", got)
	}
}

func TestThenByDescending(t *testing.T) {
	byAge := linq.OrderBy(linq.From(emps), func(e emp) int { return e.Age })
	got := linq.ThenByDescending(byAge, func(e emp) string { return e.Name }).ToSlice()
	want := []emp{{"Dee", 25}, {"Ann", 25}, {"Cy", 30}, {"Bob", 30}}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v", got)
	}
}

func TestOrderByDescendingThenBy(t *testing.T) {
	byAge := linq.OrderByDescending(linq.From(emps), func(e emp) int { return e.Age })
	got := linq.ThenBy(byAge, func(e emp) string { return e.Name }).ToSlice()
	want := []emp{{"Bob", 30}, {"Cy", 30}, {"Ann", 25}, {"Dee", 25}}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v", got)
	}
}

func TestThenByBranchesAreIndependent(t *testing.T) {
	base := linq.OrderBy(linq.From(emps), func(e emp) int { return e.Age })
	asc := linq.ThenBy(base, func(e emp) string { return e.Name })
	desc := linq.ThenByDescending(base, func(e emp) string { return e.Name })
	if got := asc.ToSlice(); got[0].Name != "Ann" {
		t.Fatalf("asc got %v", got)
	}
	if got := desc.ToSlice(); got[0].Name != "Dee" {
		t.Fatalf("desc got %v", got)
	}
}

func TestChainsAfterOrdering(t *testing.T) {
	ordered := linq.OrderBy(linq.From([]int{5, 3, 8, 1}), ident).
		Where(func(x int) bool { return x > 1 })
	got := linq.Select(ordered, func(x int) int { return x * 2 }).ToSlice()
	if !slices.Equal(got, []int{6, 10, 16}) {
		t.Fatalf("got %v", got)
	}
}

func TestSelectOnOrderedQueryField(t *testing.T) {
	ordered := linq.OrderBy(linq.From([]int{2, 1}), ident)
	got := linq.Select(ordered.Query, func(x int) int { return x + 1 }).ToSlice()
	if !slices.Equal(got, []int{2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestOrderByDoesNotMutateSource(t *testing.T) {
	data := []int{3, 1, 2}
	linq.OrderBy(linq.From(data), ident).ToSlice()
	if !slices.Equal(data, []int{3, 1, 2}) {
		t.Fatalf("source was modified: %v", data)
	}
}

func TestOrderByDeferredAndReEnumerable(t *testing.T) {
	calls := 0
	q := linq.OrderBy(linq.From([]int{3, 1, 2}), func(x int) int {
		calls++
		return x
	})
	if calls != 0 {
		t.Fatalf("key selector ran before enumeration")
	}
	if !slices.Equal(q.ToSlice(), []int{1, 2, 3}) || !slices.Equal(q.ToSlice(), []int{1, 2, 3}) {
		t.Fatal("enumerations differ")
	}
	if calls == 0 {
		t.Fatal("key selector never ran")
	}
}
