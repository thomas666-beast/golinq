# golinq

LINQ-style lazy query operators for Go, built on `iter.Seq`.

```bash
go get github.com/thomas666-beast/golinq
```

Requires Go 1.23 or newer.

## Example

```go
import linq "github.com/thomas666-beast/golinq"

people := []Person{
    {"Ann", 25, "Munich"},
    {"Bob", 30, "Berlin"},
    {"Cy", 35, "Munich"},
}

over30 := linq.From(people).Where(func(p Person) bool { return p.Age >= 30 })
sorted := linq.OrderBy(over30, func(p Person) string { return p.Name })
names := linq.Select(sorted.Query, func(p Person) string { return p.Name }).ToSlice()
// [Bob Cy]

for g := range linq.GroupBy(linq.From(people), func(p Person) string { return p.City }).Seq() {
    fmt.Println(g.Key, g.Count())
}
// Munich 2
// Berlin 1
```

## Methods vs functions

Go methods can't have their own type parameters, so:

- **Methods** (they chain): `Where`, `Take`, `Skip`, `TakeWhile`, `SkipWhile`,
  `Concat`, `Reverse`, `First`, `Last`, `Single`, `ElementAt`, `Any`, `All`,
  `Count`, `Reduce`, `ToSlice`, `Seq`
- **Functions** (they change the type or need a constraint): `Select`,
  `SelectMany`, `OrderBy`, `ThenBy`, `Chunk`, `Contains`, `Sum`, `Min`, `Max`,
  `MinBy`, `MaxBy`, `Average`, `Aggregate`, `Distinct`, `DistinctBy`, `Union`,
  `Intersect`, `Except`, `GroupBy`, `Join`, `GroupJoin`, `Zip`, `ToMap`, `ToSet`

Variants that C# expresses as overloads get their own names, such as
`FirstWhere`, `CountWhere`, `SumBy` and `GroupByElement`.

## Behavior notes

- **Lazy and re-enumerable.** Operators returning a `Query` do nothing until
  enumerated, and each enumeration re-runs the query. A query built with
  `FromSeq` on a one-shot sequence can only be consumed once.
- **Early exit.** Operators stop pulling from the source as soon as the
  consumer stops, so `Take`, `First` and `Any` work on infinite sequences.
- **No panics on empty input.** `First`, `Last`, `Min`, `Max`, `Average` and
  friends return `(value, ok)`. `Single` and `ToMap` return errors
  (`ErrNoElements`, `ErrMultipleElements`, `ErrDuplicateKey`).
  `Chunk` panics if `size < 1`.
- **Ordering.** `OrderBy` keys must be `cmp.Ordered`. The sort is stable and
  never modifies the source slice. `OrderedQuery` and `Grouping` embed
  `Query`; pass their `.Query` field to functions such as `Select` or `Sum`.
- **Comparable keys.** Set operations, `GroupBy`, `Join` and `ToMap` need
  `comparable` keys.

## Development

```bash
go test ./...
go test -bench=. -benchmem -run='^$'
```

## License

MIT
