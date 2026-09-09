# safe-sum

The program receives one input-file path. The editor contains
`safe-sum/main.go`; implement only:

```go
func sumLine(line string) string
```

`line` is one raw input line. Its tokens are separated by whitespace. Every
valid token is a signed base-10 `int64`; you may assume the sum also fits in an
`int64`. Return their sum as a decimal string. An empty or whitespace-only line
returns `"0"`.

Process tokens from left to right. At the first invalid token, return exactly:

```text
error: invalid token "<tok>"
```

Replace `<tok>` with that token's original text. Do not process later tokens on
the line, and never panic on invalid input.

File reading, line splitting, joining, and printing are provided. The driver
calls `sumLine` independently for every line, preserves blank lines, joins the
returned strings with `\n`, and prints without a final newline.

Examples for individual calls:

- `sumLine("+7 -3")` returns `"4"`.
- `sumLine("5 x 7")` returns `"error: invalid token \"x\""`.

```
flowchart LR
    F["input-file path"] --> R["read + split lines · provided"]
    R -->|"line · string"| S["sumLine · TODO"]
    S -->|"sum or first-token error · string"| J["join + print · provided"]
    J -->|"results separated by LF; no final LF"| O["stdout"]
```

## Useful references

- https://pkg.go.dev/strings#Fields
- https://pkg.go.dev/strconv#ParseInt
- https://pkg.go.dev/strconv#FormatInt
