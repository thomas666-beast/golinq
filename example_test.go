package linq_test

import (
	"errors"
	"fmt"

	linq "github.com/thomas666-beast/golinq"
)

type Person struct {
	Name string
	Age  int
	City string
}

func Example() {
	people := []Person{
		{"Ann", 25, "Munich"},
		{"Bob", 30, "Berlin"},
		{"Cy", 35, "Munich"},
	}

	over30 := linq.From(people).Where(func(p Person) bool { return p.Age >= 30 })
	sorted := linq.OrderBy(over30, func(p Person) string { return p.Name })
	names := linq.Select(sorted.Query, func(p Person) string { return p.Name })
	fmt.Println(names.ToSlice())

	byCity := linq.GroupBy(linq.From(people), func(p Person) string { return p.City })
	for g := range byCity.Seq() {
		fmt.Println(g.Key, g.Count())
	}
	// Output:
	// [Bob Cy]
	// Munich 2
	// Berlin 1
}

func ExampleQuery_Take() {
	naturals := linq.FromSeq(func(yield func(int) bool) {
		for i := 1; ; i++ {
			if !yield(i) {
				return
			}
		}
	})
	evens := naturals.Where(func(n int) bool { return n%2 == 0 })
	fmt.Println(evens.Take(3).ToSlice())
	// Output: [2 4 6]
}

func ExampleOrderBy() {
	words := linq.From([]string{"pear", "fig", "apple", "kiwi", "plum"})
	byLength := linq.OrderBy(words, func(s string) int { return len(s) })
	thenAlpha := linq.ThenBy(byLength, func(s string) string { return s })
	fmt.Println(thenAlpha.ToSlice())
	// Output: [fig kiwi pear plum apple]
}

func ExampleJoin() {
	pairs := linq.Join(
		linq.From(customers), linq.From(orders),
		func(c customer) int { return c.ID },
		func(o order) int { return o.CustomerID },
		func(c customer, o order) string { return c.Name + " bought " + o.Item },
	)
	for s := range pairs.Seq() {
		fmt.Println(s)
	}
	// Output:
	// Ann bought book
	// Ann bought lamp
	// Bob bought pen
}

func ExampleToMap() {
	_, err := linq.ToMap(
		linq.From([]string{"a", "b"}),
		func(s string) int { return len(s) },
	)
	fmt.Println(errors.Is(err, linq.ErrDuplicateKey))
	// Output: true
}
