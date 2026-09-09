# bounded-fanout

The editor contains `fanout/main.go`. Implement only:

```go
func squareAll(nums []int) []int
```

Return a slice where `out[i]` is the square of `nums[i]`. Results must remain
in input order, regardless of the order in which concurrent work finishes.

Compute the results concurrently using a bounded pool of at most `workers`
goroutines. Use the provided `jobs` channel and `wg` wait group; do not start one
goroutine per input value. The function must wait for every job and leave no
goroutine running, including when `nums` is empty.

The provided driver reads signed integers from the input-file path, calls
`squareAll`, and prints one result per line without an added final newline.

```
flowchart LR
    F["input file path"] --> P["read + parse · provided"]
    P -->|"[]int"| S["squareAll · TODO"]
    S -->|"[]int in input order"| O["print · provided"]
    O --> X["stdout"]
```

Example: input values `3, 5, 1` produce `9, 25, 1` in that order.

## Useful references

- https://pkg.go.dev/sync#WaitGroup
- https://go.dev/ref/spec#Channel_types
- https://go.dev/ref/spec#Go_statements
