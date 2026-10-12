package linq_test

import (
	"slices"
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

func TestFromSlice(t *testing.T) {
	got := linq.From([]int{1, 2, 3}).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestFromEmptyAndNil(t *testing.T) {
	if got := linq.From([]int{}).ToSlice(); got == nil || len(got) != 0 {
		t.Fatalf("expected empty non-nil slice, got %#v", got)
	}
	if got := linq.From[int](nil).ToSlice(); got == nil || len(got) != 0 {
		t.Fatalf("expected empty non-nil slice, got %#v", got)
	}
}

func TestFromSeq(t *testing.T) {
	seq := func(yield func(int) bool) {
		for i := 0; i < 3; i++ {
			if !yield(i) {
				return
			}
		}
	}
	got := linq.FromSeq(seq).ToSlice()
	if !slices.Equal(got, []int{0, 1, 2}) {
		t.Fatalf("got %v", got)
	}
}

func TestSeqWorksWithRange(t *testing.T) {
	var got []string
	for s := range linq.From([]string{"a", "b", "c"}).Seq() {
		got = append(got, s)
	}
	if !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Fatalf("got %v", got)
	}
}

func TestSeqSupportsEarlyBreak(t *testing.T) {
	infinite := func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
	var got []int
	for x := range linq.FromSeq(infinite).Seq() {
		if x == 3 {
			break
		}
		got = append(got, x)
	}
	if !slices.Equal(got, []int{0, 1, 2}) {
		t.Fatalf("got %v", got)
	}
}

func TestQueryIsReEnumerable(t *testing.T) {
	q := linq.From([]int{1, 2, 3})
	if !slices.Equal(q.ToSlice(), []int{1, 2, 3}) || !slices.Equal(q.ToSlice(), []int{1, 2, 3}) {
		t.Fatal("second enumeration differs")
	}
}

func TestQueryReflectsSourceChanges(t *testing.T) {
	data := []int{1, 2}
	q := linq.From(data)
	data[0] = 99
	if got := q.ToSlice(); !slices.Equal(got, []int{99, 2}) {
		t.Fatalf("got %v", got)
	}
}

func TestZeroValueQueryIsEmpty(t *testing.T) {
	var q linq.Query[int]
	if got := q.ToSlice(); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}
