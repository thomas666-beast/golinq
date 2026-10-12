package linq_test

import (
	"slices"
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

type customer struct {
	ID   int
	Name string
}

type order struct {
	CustomerID int
	Item       string
}

var customers = []customer{{1, "Ann"}, {2, "Bob"}, {3, "Cy"}}
var orders = []order{{1, "book"}, {2, "pen"}, {1, "lamp"}, {9, "ghost"}}

// GroupBy
func TestGroupByKey(t *testing.T) {
	groups := linq.GroupBy(linq.From([]int{1, 2, 3, 4, 5, 6}), func(x int) int { return x % 2 }).ToSlice()
	if len(groups) != 2 {
		t.Fatalf("got %d groups", len(groups))
	}
	if groups[0].Key != 1 || !slices.Equal(groups[0].ToSlice(), []int{1, 3, 5}) {
		t.Fatalf("group 0: %v %v", groups[0].Key, groups[0].ToSlice())
	}
	if groups[1].Key != 0 || !slices.Equal(groups[1].ToSlice(), []int{2, 4, 6}) {
		t.Fatalf("group 1: %v %v", groups[1].Key, groups[1].ToSlice())
	}
}

func TestGroupByFirstAppearanceOrder(t *testing.T) {
	groups := linq.GroupBy(linq.From([]string{"b1", "a1", "b2", "a2"}), func(s string) byte { return s[0] }).ToSlice()
	if groups[0].Key != 'b' || groups[1].Key != 'a' {
		t.Fatalf("got keys %c %c", groups[0].Key, groups[1].Key)
	}
}

func TestGroupByElement(t *testing.T) {
	groups := linq.GroupByElement(
		linq.From([]string{"apple", "avocado", "banana"}),
		func(w string) byte { return w[0] },
		func(w string) int { return len(w) },
	).ToSlice()
	if groups[0].Key != 'a' || !slices.Equal(groups[0].ToSlice(), []int{5, 7}) {
		t.Fatalf("group 0: %v", groups[0].ToSlice())
	}
	if groups[1].Key != 'b' || !slices.Equal(groups[1].ToSlice(), []int{6}) {
		t.Fatalf("group 1: %v", groups[1].ToSlice())
	}
}

func TestGroupByResultViaSelect(t *testing.T) {
	groups := linq.GroupBy(linq.From([]int{1, 2, 3, 4, 5}), func(x int) int { return x % 2 })
	type kv struct{ Key, Sum int }
	got := linq.Select(groups, func(g linq.Grouping[int, int]) kv {
		return kv{g.Key, linq.Sum(g.Query)}
	}).ToSlice()
	if !slices.Equal(got, []kv{{1, 9}, {0, 6}}) {
		t.Fatalf("got %v", got)
	}
}

func TestGroupByEmpty(t *testing.T) {
	got := linq.GroupBy(linq.From([]int{}), func(x int) int { return x }).ToSlice()
	if len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestGroupIsAQuery(t *testing.T) {
	g, _ := linq.GroupBy(linq.From([]int{1, 2, 3, 4}), func(x int) int { return x % 2 }).First()
	if got := g.Where(func(x int) bool { return x > 1 }).ToSlice(); !slices.Equal(got, []int{3}) {
		t.Fatalf("got %v", got)
	}
	if g.Count() != 2 {
		t.Fatalf("count %d", g.Count())
	}
}

func TestGroupByDeferredExecution(t *testing.T) {
	var calls []int
	q := linq.GroupBy(linq.From([]int{1, 2, 3}), func(x int) int {
		calls = append(calls, x)
		return x % 2
	})
	if len(calls) != 0 {
		t.Fatalf("key ran before enumeration: %v", calls)
	}
	q.ToSlice()
	if !slices.Equal(calls, []int{1, 2, 3}) {
		t.Fatalf("calls %v", calls)
	}
}

func TestGroupByIsReEnumerable(t *testing.T) {
	q := linq.GroupBy(linq.From([]int{1, 2, 3}), func(x int) int { return x % 2 })
	if len(q.ToSlice()) != 2 || len(q.ToSlice()) != 2 {
		t.Fatal("enumerations differ")
	}
}

// Join
func TestJoin(t *testing.T) {
	got := linq.Join(
		linq.From(customers), linq.From(orders),
		func(c customer) int { return c.ID },
		func(o order) int { return o.CustomerID },
		func(c customer, o order) string { return c.Name + ":" + o.Item },
	).ToSlice()
	if !slices.Equal(got, []string{"Ann:book", "Ann:lamp", "Bob:pen"}) {
		t.Fatalf("got %v", got)
	}
}

func TestJoinExcludesUnmatched(t *testing.T) {
	id := func(x int) int { return x }
	got := linq.Join(
		linq.From([]int{1, 2, 3}), linq.From([]int{2, 3, 4}), id, id,
		func(a, b int) [2]int { return [2]int{a, b} },
	).ToSlice()
	if !slices.Equal(got, [][2]int{{2, 2}, {3, 3}}) {
		t.Fatalf("got %v", got)
	}
}

func TestJoinMultipleMatchesOnBothSides(t *testing.T) {
	got := linq.Join(
		linq.From([]int{1, 1}), linq.From([]string{"a", "b"}),
		func(x int) int { return 1 },
		func(s string) int { return 1 },
		func(x int, s string) string { return s },
	).ToSlice()
	if !slices.Equal(got, []string{"a", "b", "a", "b"}) {
		t.Fatalf("got %v", got)
	}
}

func TestJoinDeferredExecution(t *testing.T) {
	var calls []int
	q := linq.Join(
		linq.From([]int{1}), linq.From([]int{1, 2}),
		func(x int) int { return x },
		func(y int) int { calls = append(calls, y); return y },
		func(x, y int) int { return y },
	)
	if len(calls) != 0 {
		t.Fatalf("ran before enumeration: %v", calls)
	}
	if got := q.ToSlice(); !slices.Equal(got, []int{1}) {
		t.Fatalf("got %v", got)
	}
	if !slices.Equal(calls, []int{1, 2}) {
		t.Fatalf("calls %v", calls)
	}
}

// GroupJoin
func TestGroupJoin(t *testing.T) {
	type res struct {
		Name  string
		Items []string
	}
	got := linq.GroupJoin(
		linq.From(customers), linq.From(orders),
		func(c customer) int { return c.ID },
		func(o order) int { return o.CustomerID },
		func(c customer, os linq.Query[order]) res {
			return res{c.Name, linq.Select(os, func(o order) string { return o.Item }).ToSlice()}
		},
	).ToSlice()
	if len(got) != 3 ||
		!slices.Equal(got[0].Items, []string{"book", "lamp"}) ||
		!slices.Equal(got[1].Items, []string{"pen"}) ||
		len(got[2].Items) != 0 || got[2].Name != "Cy" {
		t.Fatalf("got %v", got)
	}
}

func TestGroupJoinIncludesOuterWithoutMatches(t *testing.T) {
	id := func(x int) int { return x }
	got := linq.GroupJoin(
		linq.From([]int{1, 2}), linq.From([]int{}), id, id,
		func(x int, ys linq.Query[int]) [2]int { return [2]int{x, ys.Count()} },
	).ToSlice()
	if !slices.Equal(got, [][2]int{{1, 0}, {2, 0}}) {
		t.Fatalf("got %v", got)
	}
}

func TestGroupJoinInnerCanBeSeq(t *testing.T) {
	id := func(x int) int { return x }
	inner := linq.FromSeq(slices.Values([]int{2, 2, 3}))
	got := linq.GroupJoin(
		linq.From([]int{1, 2}), inner, id, id,
		func(x int, ys linq.Query[int]) [2]int { return [2]int{x, ys.Count()} },
	).ToSlice()
	if !slices.Equal(got, [][2]int{{1, 0}, {2, 2}}) {
		t.Fatalf("got %v", got)
	}
}
