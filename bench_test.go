package linq_test

import (
	"testing"

	linq "github.com/thomas666-beast/golinq"
)

func benchData() []int {
	data := make([]int, 1_000_000)
	for i := range data {
		data[i] = i
	}
	return data
}

var sink int

func BenchmarkPlainLoop(b *testing.B) {
	data := benchData()
	b.ReportAllocs()
	for b.Loop() {
		total := 0
		for _, x := range data {
			if x%2 == 0 {
				total += x * 2
			}
		}
		sink = total
	}
}

func BenchmarkWhereSelectSum(b *testing.B) {
	data := benchData()
	b.ReportAllocs()
	for b.Loop() {
		q := linq.Select(
			linq.From(data).Where(func(x int) bool { return x%2 == 0 }),
			func(x int) int { return x * 2 },
		)
		sink = linq.Sum(q)
	}
}

func BenchmarkOrderBy(b *testing.B) {
	data := benchData()
	b.ReportAllocs()
	for b.Loop() {
		sink = len(linq.OrderByDescending(linq.From(data), ident).ToSlice())
	}
}

func BenchmarkGroupBy(b *testing.B) {
	data := benchData()
	b.ReportAllocs()
	for b.Loop() {
		sink = linq.GroupBy(linq.From(data), func(x int) int { return x % 100 }).Count()
	}
}
