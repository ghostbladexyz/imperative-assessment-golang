# tetris

The editor contains `tetris/main.go`. Implement these two functions:

```go
func canPlace(grid [][]byte, p Piece, row, col int) bool
func stamp(grid [][]byte, p Piece, row, col int, ch byte)
```

`Piece` is a slice of normalized `Cell{R, C}` coordinates. The provided parser
converts `#` cells from blank-line-separated shapes into `Piece` values and
moves each piece's top-left occupied cell to `(0, 0)`. A grid is an `n × n`
slice whose empty cells contain `'.'`; `(row, col)` is a candidate offset.

`canPlace` returns `true` only when every translated coordinate
`(row + cell.R, col + cell.C)` is inside the grid and currently `'.'`. It must
not mutate the grid. Return `false` as soon as either condition is not met.

`stamp` may assume the placement fits. For every cell in `p`, write `ch` at its
translated grid coordinate. The provided search passes a letter byte to place a
piece and `'.'` to undo that placement.

Inputs contain at most 4 valid tetrominoes. Pieces may be translated only—no
rotation—and may not overlap. Backtracking, minimum-square search, and printing
are also provided; you do not need to redesign them.

The provided printer writes each grid row followed by a newline. Your two
functions must produce no other output.

```
flowchart LR
    F["file path · string"] --> P["parse pieces · provided"]
    P -->|"[]Piece"| B["backtrack + minimum side · provided"]
    B -->|"grid [][]byte; p Piece; row, col int"| C["canPlace · TODO<br/>returns bool"]
    B -->|"grid, p, row, col; ch byte"| S["stamp · TODO<br/>mutates grid"]
    B -->|"[][]byte"| O["print rows · provided"]
    O -->|"one row per line"| X["stdout"]
```

For example, if `p` contains cells `(0,0)` and `(0,1)`, offset `(2,3)` refers
to `grid[2][3]` and `grid[2][4]`.

## Useful references

- https://go.dev/ref/spec#Slice_types
- https://go.dev/ref/spec#Struct_types
- https://go.dev/ref/spec#Index_expressions
