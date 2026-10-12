package linq_test

import (
	"strconv"
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

type celsius float64

// Count
func TestCount(t *testing.T) {
	if got := linq.From([]int{1, 2, 3}).Count(); got != 3 {
		t.Fatalf("got %d", got)
	}
}

func TestCountEmpty(t *testing.T) {
	if got := linq.From([]int{}).Count(); got != 0 {
		t.Fatalf("got %d", got)
	}
}

func TestCountWhere(t *testing.T) {
	got := linq.From([]int{1, 2, 3, 4}).CountWhere(func(x int) bool { return x%2 == 0 })
	if got != 2 {
		t.Fatalf("got %d", got)
	}
}

// Sum
func TestSum(t *testing.T) {
	if got := linq.Sum(linq.From([]int{1, 2, 3})); got != 6 {
		t.Fatalf("got %d", got)
	}
}

func TestSumEmpty(t *testing.T) {
	if got := linq.Sum(linq.From([]int{})); got != 0 {
		t.Fatalf("got %d", got)
	}
}

func TestSumBy(t *testing.T) {
	got := linq.SumBy(linq.From([]string{"a", "bb", "ccc"}), func(s string) int { return len(s) })
	if got != 6 {
		t.Fatalf("got %d", got)
	}
}

func TestSumFloats(t *testing.T) {
	if got := linq.Sum(linq.From([]float64{0.5, 1.5})); got != 2.0 {
		t.Fatalf("got %v", got)
	}
}

func TestSumNamedType(t *testing.T) {
	got := linq.Sum(linq.From([]celsius{10.5, 20}))
	if got != celsius(30.5) {
		t.Fatalf("got %v", got)
	}
}

// Min / Max
func TestMin(t *testing.T) {
	v, ok := linq.Min(linq.From([]int{3, 1, 2}))
	if !ok || v != 1 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestMax(t *testing.T) {
	v, ok := linq.Max(linq.From([]int{3, 1, 2}))
	if !ok || v != 3 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestMinMaxStrings(t *testing.T) {
	q := linq.From([]string{"pear", "apple", "zebra"})
	if v, _ := linq.Min(q); v != "apple" {
		t.Fatalf("min got %v", v)
	}
	if v, _ := linq.Max(q); v != "zebra" {
		t.Fatalf("max got %v", v)
	}
}

func TestMinMaxSingleElement(t *testing.T) {
	q := linq.From([]int{5})
	if v, ok := linq.Min(q); !ok || v != 5 {
		t.Fatalf("min got %v, %v", v, ok)
	}
	if v, ok := linq.Max(q); !ok || v != 5 {
		t.Fatalf("max got %v, %v", v, ok)
	}
}

func TestMinEmpty(t *testing.T) {
	if _, ok := linq.Min(linq.From([]int{})); ok {
		t.Fatal("expected not ok")
	}
}

func TestMaxEmpty(t *testing.T) {
	if _, ok := linq.Max(linq.From([]int{})); ok {
		t.Fatal("expected not ok")
	}
}

// MinBy / MaxBy
func TestMinByMaxBy(t *testing.T) {
	q := linq.From([]string{"ccc", "a", "bb"})
	length := func(s string) int { return len(s) }
	if v, ok := linq.MinBy(q, length); !ok || v != "a" {
		t.Fatalf("minBy got %v, %v", v, ok)
	}
	if v, ok := linq.MaxBy(q, length); !ok || v != "ccc" {
		t.Fatalf("maxBy got %v, %v", v, ok)
	}
}

func TestMinByMaxByTiesPickFirst(t *testing.T) {
	age := func(e emp) int { return e.Age }
	if v, _ := linq.MinBy(linq.From(emps), age); v.Name != "Ann" {
		t.Fatalf("minBy got %v", v)
	}
	if v, _ := linq.MaxBy(linq.From(emps), age); v.Name != "Bob" {
		t.Fatalf("maxBy got %v", v)
	}
}

func TestMinByEmpty(t *testing.T) {
	_, ok := linq.MinBy(linq.From([]emp{}), func(e emp) int { return e.Age })
	if ok {
		t.Fatal("expected not ok")
	}
}

func TestMaxByEmpty(t *testing.T) {
	_, ok := linq.MaxBy(linq.From([]emp{}), func(e emp) int { return e.Age })
	if ok {
		t.Fatal("expected not ok")
	}
}

// Average
func TestAverage(t *testing.T) {
	v, ok := linq.Average(linq.From([]int{1, 2, 3, 4}))
	if !ok || v != 2.5 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestAverageBy(t *testing.T) {
	v, ok := linq.AverageBy(linq.From([]string{"a", "bbb"}), func(s string) int { return len(s) })
	if !ok || v != 2.0 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestAverageEmpty(t *testing.T) {
	if _, ok := linq.Average(linq.From([]int{})); ok {
		t.Fatal("expected not ok")
	}
}

// Aggregate / Reduce
func TestAggregateWithSeed(t *testing.T) {
	got := linq.Aggregate(linq.From([]int{1, 2, 3}), 100, func(acc, x int) int { return acc + x })
	if got != 106 {
		t.Fatalf("got %d", got)
	}
}

func TestAggregateChangesType(t *testing.T) {
	got := linq.Aggregate(linq.From([]int{1, 2, 3}), "", func(acc string, x int) string {
		return acc + strconv.Itoa(x)
	})
	if got != "123" {
		t.Fatalf("got %q", got)
	}
}

func TestAggregateBuildsMap(t *testing.T) {
	got := linq.Aggregate(linq.From([]string{"a", "bb"}), map[string]int{}, func(acc map[string]int, s string) map[string]int {
		acc[s] = len(s)
		return acc
	})
	if got["a"] != 1 || got["bb"] != 2 {
		t.Fatalf("got %v", got)
	}
}

func TestAggregateEmptyReturnsSeed(t *testing.T) {
	got := linq.Aggregate(linq.From([]int{}), 42, func(acc, x int) int { return acc + x })
	if got != 42 {
		t.Fatalf("got %d", got)
	}
}

func TestReduce(t *testing.T) {
	v, ok := linq.From([]int{1, 2, 3, 4}).Reduce(func(acc, x int) int { return acc + x })
	if !ok || v != 10 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestReduceEmpty(t *testing.T) {
	if _, ok := linq.From([]int{}).Reduce(func(acc, x int) int { return acc + x }); ok {
		t.Fatal("expected not ok")
	}
}

func TestReduceSingleElementNeverCallsFunc(t *testing.T) {
	v, ok := linq.From([]int{7}).Reduce(func(acc, x int) int {
		t.Fatal("fn should not be called")
		return 0
	})
	if !ok || v != 7 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestAggregationChainsAfterWhere(t *testing.T) {
	evens := linq.From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
		Where(func(x int) bool { return x%2 == 0 })
	if got := linq.Sum(evens); got != 30 {
		t.Fatalf("got %d", got)
	}
}
