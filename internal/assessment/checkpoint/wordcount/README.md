# wordcount

Complete `orderedWords` in `wordcount/main.go`.

The provided program receives one input-file path, counts whitespace-separated
words, and prints one `word count` line for each distinct word. Your function
must return the words in this order:

1. higher count first;
2. for equal counts, word ascending.

The order must be deterministic. Go map iteration order is not stable.

The provided formatter joins output lines with `\n` and does not add a final
newline. Punctuation is part of a word: `go` and `go,` are different words.

Example input:

```text
go rust go java rust go java rust
```

Exact output (there is no newline after `java 2`):

```text
go 3
rust 3
java 2
```

The file reading, counting, and output formatting are provided:

```
flowchart LR
    A["input file"] --> B["provided counting"]
    B --> C["orderedWords — TODO"]
    C --> D["provided formatting"]
    D --> E["stdout"]
```

## Useful references

- https://pkg.go.dev/sort#Slice
- https://go.dev/blog/maps
- https://go.dev/ref/spec#For_range
