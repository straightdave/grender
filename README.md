Grender
=======

[![Go Reference](https://pkg.go.dev/badge/github.com/straightdave/grender.svg)](https://pkg.go.dev/github.com/straightdave/grender)

Yet another rendering tool for Golang.

## Example

### Using layout

```go
r := NewGrender()
r.AddLayout("L1", `<layout>{{ yield }}</layout>`)
r.Add("P1", `Hello {{ .name }}`)

out, _ := r.Render("L1", "P1", map[string]interface{}{
    "name": "dave",
})

// out => "<layout>dave</layout>"
```

### Using shared templates

```go
r := NewGrender()
r.Add("S1", `<shared>{{ .name }}</shared>`)
r.Add("P1", `Any Template Can Use {{ share "S1" }}`)

// this case doesn't use any layout
out, _ := r.Render("", "P1", map[string]interface{}{
    "name": "dave",
})
// out => "Any Template Can Use <shared>dave</shared>"
```

Use shared templates with both layout and page:

```go
r := NewGrender()
r.AddLayout("L1", `<layout>{{ share "S1" }} -> {{ yield }}</layout>`)
r.Add("S1", `<shared>{{ .name }}</shared>`)
r.Add("P1", `Any Template Can Use {{ share "S1" }}`)

out, _ := r.Render("L1", "P1", map[string]interface{}{
    "name": "dave",
})
// out => "<layout><shared>dave</shared> -> Any Template Can Use <shared>dave</shared></layout>"
```
