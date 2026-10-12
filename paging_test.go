package linq_test

import (
	"slices"
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

func infiniteFrom1() linq.Query[int] {
	return linq.FromSeq(func(yield func(int) bool) {
		for i := 1; ; i++ {
			if !yield(i) {
				return
			}
		}
	})
}

func TestTake(t *testing.T) {
	got := linq.From([]int{1, 2, 3, 4, 5}).Take(3).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestTakeMoreThanAvailable(t *testing.T) {
	got := linq.From([]int{1, 2}).Take(10).ToSlice()
	if !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("got %v", got)
	}
}

func TestTakeZeroAndNegative(t *testing.T) {
	if got := linq.From([]int{1, 2, 3}).Take(0).ToSlice(); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
	if got := linq.From([]int{1, 2, 3}).Take(-5).ToSlice(); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestTakeWorksOnInfiniteSequence(t *testing.T) {
	got := infiniteFrom1().Take(3).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestTakeDoesNotPullExtra(t *testing.T) {
	var pulled []int
	src := linq.FromSeq(func(yield func(int) bool) {
		for i := 0; i < 100; i++ {
			pulled = append(pulled, i)
			if !yield(i) {
				return
			}
		}
	})
	src.Take(3).ToSlice()
	if !slices.Equal(pulled, []int{0, 1, 2}) {
		t.Fatalf("pulled %v", pulled)
	}
}

func TestSkip(t *testing.T) {
	got := linq.From([]int{1, 2, 3, 4, 5}).Skip(2).ToSlice()
	if !slices.Equal(got, []int{3, 4, 5}) {
		t.Fatalf("got %v", got)
	}
}

func TestSkipMoreThanAvailable(t *testing.T) {
	if got := linq.From([]int{1, 2}).Skip(10).ToSlice(); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestSkipZeroAndNegative(t *testing.T) {
	if got := linq.From([]int{1, 2, 3}).Skip(0).ToSlice(); !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
	if got := linq.From([]int{1, 2, 3}).Skip(-5).ToSlice(); !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestSkipTakePaging(t *testing.T) {
	nums := make([]int, 20)
	for i := range nums {
		nums[i] = i + 1
	}
	got := linq.From(nums).Skip(10).Take(5).ToSlice()
	if !slices.Equal(got, []int{11, 12, 13, 14, 15}) {
		t.Fatalf("got %v", got)
	}
}

func TestTakeWhile(t *testing.T) {
	got := linq.From([]int{1, 2, 3, 10, 1, 2}).
		TakeWhile(func(x int) bool { return x < 5 }).ToSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("got %v", got)
	}
}

func TestTakeWhileNoneMatch(t *testing.T) {
	got := linq.From([]int{5, 1}).TakeWhile(func(x int) bool { return x < 5 }).ToSlice()
	if len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestSkipWhile(t *testing.T) {
	got := linq.From([]int{1, 2, 3, 10, 1, 2}).
		SkipWhile(func(x int) bool { return x < 5 }).ToSlice()
	if !slices.Equal(got, []int{10, 1, 2}) {
		t.Fatalf("got %v", got)
	}
}

func TestSkipWhileAllMatch(t *testing.T) {
	got := linq.From([]int{1, 2, 3}).SkipWhile(func(x int) bool { return x < 5 }).ToSlice()
	if len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestChunk(t *testing.T) {
	got := linq.Chunk(linq.From([]int{1, 2, 3, 4, 5}), 2).ToSlice()
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Fatalf("got %v", got)
		}
	}
}

func TestChunkExactMultiple(t *testing.T) {
	got := linq.Chunk(linq.From([]int{1, 2, 3, 4}), 2).ToSlice()
	if len(got) != 2 || !slices.Equal(got[1], []int{3, 4}) {
		t.Fatalf("got %v", got)
	}
}

func TestChunkEmpty(t *testing.T) {
	if got := linq.Chunk(linq.From([]int{}), 3).ToSlice(); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestChunkInvalidSizePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	linq.Chunk(linq.From([]int{1, 2, 3}), 0)
}

func TestChunkSlicesAreIndependent(t *testing.T) {
	got := linq.Chunk(linq.From([]int{1, 2, 3, 4}), 2).ToSlice()
	got[0][0] = 99
	if got[1][0] != 3 {
		t.Fatalf("chunks share memory: %v", got)
	}
}

func TestChunkStopsOnEarlyBreak(t *testing.T) {
	var pulled []int
	src := linq.FromSeq(func(yield func(int) bool) {
		for i := 0; i < 100; i++ {
			pulled = append(pulled, i)
			if !yield(i) {
				return
			}
		}
	})
	for range linq.Chunk(src, 2).Seq() {
		break
	}
	if !slices.Equal(pulled, []int{0, 1}) {
		t.Fatalf("pulled %v", pulled)
	}
}
