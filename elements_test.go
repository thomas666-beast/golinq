package linq_test

import (
	"errors"
	"slices"
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

func trackedSource(pulled *[]int) linq.Query[int] {
	return linq.FromSeq(func(yield func(int) bool) {
		for i := 0; i < 100; i++ {
			*pulled = append(*pulled, i)
			if !yield(i) {
				return
			}
		}
	})
}

func TestFirst(t *testing.T) {
	v, ok := linq.From([]int{1, 2, 3}).First()
	if !ok || v != 1 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestFirstEmpty(t *testing.T) {
	v, ok := linq.From([]int{}).First()
	if ok || v != 0 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestFirstWhere(t *testing.T) {
	v, ok := linq.From([]int{1, 2, 3, 4}).FirstWhere(func(x int) bool { return x > 2 })
	if !ok || v != 3 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestFirstWhereNoMatch(t *testing.T) {
	_, ok := linq.From([]int{1, 2}).FirstWhere(func(x int) bool { return x > 10 })
	if ok {
		t.Fatal("expected no match")
	}
}

func TestFirstStopsEarly(t *testing.T) {
	var pulled []int
	trackedSource(&pulled).First()
	if !slices.Equal(pulled, []int{0}) {
		t.Fatalf("pulled %v", pulled)
	}
}

func TestSingle(t *testing.T) {
	v, err := linq.From([]int{7}).Single()
	if err != nil || v != 7 {
		t.Fatalf("got %v, %v", v, err)
	}
}

func TestSingleEmpty(t *testing.T) {
	_, err := linq.From([]int{}).Single()
	if !errors.Is(err, linq.ErrNoElements) {
		t.Fatalf("got %v", err)
	}
}

func TestSingleMultiple(t *testing.T) {
	_, err := linq.From([]int{1, 2}).Single()
	if !errors.Is(err, linq.ErrMultipleElements) {
		t.Fatalf("got %v", err)
	}
}

func TestSingleWhere(t *testing.T) {
	v, err := linq.From([]int{1, 2, 3}).SingleWhere(func(x int) bool { return x > 2 })
	if err != nil || v != 3 {
		t.Fatalf("got %v, %v", v, err)
	}
}

func TestSingleStopsAtSecondMatch(t *testing.T) {
	var pulled []int
	_, err := trackedSource(&pulled).Single()
	if !errors.Is(err, linq.ErrMultipleElements) {
		t.Fatalf("got %v", err)
	}
	if !slices.Equal(pulled, []int{0, 1}) {
		t.Fatalf("pulled %v", pulled)
	}
}

func TestLast(t *testing.T) {
	v, ok := linq.From([]int{1, 2, 3}).Last()
	if !ok || v != 3 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestLastEmpty(t *testing.T) {
	if _, ok := linq.From([]int{}).Last(); ok {
		t.Fatal("expected not ok")
	}
}

func TestLastWhere(t *testing.T) {
	v, ok := linq.From([]int{1, 2, 3, 4}).LastWhere(func(x int) bool { return x < 3 })
	if !ok || v != 2 {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestElementAt(t *testing.T) {
	v, ok := linq.From([]string{"a", "b", "c"}).ElementAt(1)
	if !ok || v != "b" {
		t.Fatalf("got %v, %v", v, ok)
	}
}

func TestElementAtOutOfRange(t *testing.T) {
	q := linq.From([]int{1, 2})
	if _, ok := q.ElementAt(5); ok {
		t.Fatal("expected not ok for index 5")
	}
	if _, ok := q.ElementAt(-1); ok {
		t.Fatal("expected not ok for index -1")
	}
}

func TestElementAtStopsEarly(t *testing.T) {
	var pulled []int
	trackedSource(&pulled).ElementAt(2)
	if !slices.Equal(pulled, []int{0, 1, 2}) {
		t.Fatalf("pulled %v", pulled)
	}
}

func TestAny(t *testing.T) {
	if !linq.From([]int{1}).Any() {
		t.Fatal("expected true")
	}
	if linq.From([]int{}).Any() {
		t.Fatal("expected false")
	}
}

func TestAnyWhere(t *testing.T) {
	q := linq.From([]int{1, 2, 3})
	if !q.AnyWhere(func(x int) bool { return x > 2 }) {
		t.Fatal("expected true")
	}
	if q.AnyWhere(func(x int) bool { return x > 5 }) {
		t.Fatal("expected false")
	}
}

func TestAnyShortCircuits(t *testing.T) {
	var pulled []int
	ok := trackedSource(&pulled).AnyWhere(func(x int) bool { return x == 2 })
	if !ok || !slices.Equal(pulled, []int{0, 1, 2}) {
		t.Fatalf("ok=%v pulled=%v", ok, pulled)
	}
}

func TestAll(t *testing.T) {
	even := func(x int) bool { return x%2 == 0 }
	if !linq.From([]int{2, 4, 6}).All(even) {
		t.Fatal("expected true")
	}
	if linq.From([]int{2, 3, 6}).All(even) {
		t.Fatal("expected false")
	}
}

func TestAllEmptyIsTrue(t *testing.T) {
	if !linq.From([]int{}).All(func(int) bool { return false }) {
		t.Fatal("expected true for empty query")
	}
}

func TestContains(t *testing.T) {
	if !linq.Contains(linq.From([]int{1, 2, 3}), 2) {
		t.Fatal("expected true")
	}
}

func TestContainsMissing(t *testing.T) {
	if linq.Contains(linq.From([]int{1, 2, 3}), 9) {
		t.Fatal("expected false")
	}
}
