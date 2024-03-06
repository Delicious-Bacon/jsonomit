# jsonomit
Package jsonomit provides JSON marshal functions to omit empty structs
and null fields. By default, the functions omit empty structs.

Provided functions can omit zero value time.Time fields, or null fields that
result from custom MarshalJSON implementations.

`go get github.com/Delicious-Bacon/jsonomit`
