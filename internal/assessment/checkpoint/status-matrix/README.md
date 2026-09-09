# status-matrix

The editor contains `server/main.go`. Implement only `rootHandler` and
`generateHandler`; the form, mux and listener are provided.

Required outcomes:

| Request | Result |
| --- | --- |
| `GET /` | status `200`, body containing the provided `<form>` |
| `GET` on any unknown path | status `404` |
| `POST /generate` with non-empty `text` | status `200`, body exactly the uppercase text with no newline |
| `POST /generate` with missing or empty `text` | status `400` |
| `POST /generate` with `text=panic` | status `500`, without panicking out of the handler |

The exact text `panic` represents an internal failure supplied by the assessment.
You may return `500` directly for that signal; you do not need to create a real
panic. Other method combinations are not assessed.

After the `400`, `404`, and `500` cases, another `GET /` must still return `200`.
Error-response bodies are not graded.

```
flowchart LR
    Q["HTTP request"] --> M["ServeMux · provided"]
    M --> R["rootHandler · TODO"]
    M --> G["generateHandler · TODO"]
    F["form constant · provided"] --> R
    R --> O["status + body"]
    G --> O
```

## Useful references

- https://pkg.go.dev/net/http#Handler
- https://pkg.go.dev/net/http#ResponseWriter
- https://pkg.go.dev/strings#ToUpper
