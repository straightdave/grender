Grender
=======

[![Go Reference](https://pkg.go.dev/badge/github.com/straightdave/grender.svg)](https://pkg.go.dev/github.com/straightdave/grender)

Yet another rendering tool for Golang.

## Example

### Render with a layout

Layouts are special templates that would be rendered in a slightly different way.

```go
import grender
r := grender.New()
r.AddLayout("L1", `<layout>{{ yield }}</layout>`)
r.Add("P1", `Hello {{ .name }}`)

out, _ := r.Render("L1", "P1", map[string]interface{}{
    "name": "dave",
})
// out => "<layout>dave</layout>"
```

### Render with shared templates

Shared templates are just templates.

```go
import grender
r := grender.New()
r.Add("S1", `<shared>{{ .name }}</shared>`)
r.Add("P1", `Any Template Can Use {{ share "S1" }}`)

// no layout used
out, _ := r.Render("", "P1", map[string]interface{}{
    "name": "dave",
})
// out => "Any Template Can Use <shared>dave</shared>"
```

Use shared templates with both layout and page:

```go
import grender
r := grender.New()
r.AddLayout("L1", `<layout>{{ share "S1" }} -> {{ yield }}</layout>`)
r.Add("S1", `<shared>{{ .name }}</shared>`)
r.Add("P1", `Any Template Can Use {{ share "S1" }}`)

out, _ := r.Render("L1", "P1", map[string]interface{}{
    "name": "dave",
})
// out => "<layout><shared>dave</shared> -> Any Template Can Use <shared>dave</shared></layout>"
```

### Load templates from FileSystem

```go
import grender
r := grender.New()
// Anything that implements io/fs.FS interface, normally an embedded one.
// about FS: https://pkg.go.dev/io/fs
// about embed: https://pkg.go.dev/embed
err := r.LoadFromFS(yourFS)
```

### Change default options

```go
import grender
r := grender.New(
    OptionMissingKeyZero(false),
    OptionTemplateDir("templates"),
    OptionLayoutDir("templates/layouts"),
    OptionTemplateExt([]string{".tmpl", ".html"}),
)
// all values above are default.
```
