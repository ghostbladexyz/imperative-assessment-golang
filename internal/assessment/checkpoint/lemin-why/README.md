# lemin-why

Complete the `minTurns` TODO in `lemin-why/main.go`. The parser, path discovery,
file handling, and printing are provided.

## Contract

The program receives one colony file path. The provided code reduces that file
to:

- `paths []int`: the positive edge count of each independent start-to-end path;
- `ants int`: the number of ants that must reach the end.

Return the fewest turns needed when ants may be divided among the paths. Only
one ant can enter the same path per turn, while different paths run in parallel.
If `x` ants use a path of `p` edges, the last one arrives on turn `p + x - 1`.
All ants must arrive. Return `0` when there are zero ants or no paths.

The provided `main` prints one base-10 integer followed by a newline. Do not add
other output.

Examples:

- `paths = [3, 4]`, `ants = 6` → `6`
- `paths = [3]`, `ants = 3` → `5`

## Skeleton map

```
flowchart LR
    F["colony file · path"] --> P["parseColony<br/>(provided)"]
    P -->|"*Colony"| D["derivePaths<br/>(provided)"]
    D -->|"[]int · edge counts"| M["minTurns<br/>(TODO)"]
    P -->|"c.Ants · int"| M
    M -->|"int · turns"| O["fmt.Println<br/>(provided)"]
```

## Useful references

- https://go.dev/ref/spec#Slice_types
- https://go.dev/ref/spec#For_statements
- https://go.dev/ref/spec#Arithmetic_operators
