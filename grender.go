// package grender is a smart templating tool.
package grender

import (
	"fmt"
	"sync"
	"text/template"
)

// Grender is the object to handle the rendering work.
type Grender struct {
	sync.RWMutex
	m              map[string]*template.Template
	missingKeyZero bool
}

// New creates one Grender object.
func New(opts ...Option) *Grender {
	r := &Grender{
		m: make(map[string]*template.Template),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// AddLayout adds a layout template.
func (r *Grender) AddLayout(name, content string) error {
	key := "layout:" + name
	t := template.New(key).Funcs(template.FuncMap{
		"yield": func() string {
			return "{{.}}"
		},
		"current": func() string {
			return name
		},
		"share": func(sharedTempleteName string) string {
			return fmt.Sprintf("{{ share \"%s\" }}", sharedTempleteName)
		},
	})

	if r.missingKeyZero {
		t.Option("missingkey=zero")
	}

	_, err := t.Parse(content)
	if err != nil {
		return err
	}

	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[key]; ok {
		return fmt.Errorf("template %s already exists", name)
	}

	r.m[key] = t
	return nil
}

// Add adds a page template.
func (r *Grender) Add(name, content string) error {
	key := "page:" + name
	t := template.New(key).Option("missingkey=zero").Funcs(template.FuncMap{
		"current": func() string {
			return name
		},
		"share": func(sharedTempleteName string) string {
			return fmt.Sprintf("{{ share \"%s\" }}", sharedTempleteName)
		},
	})

	if r.missingKeyZero {
		t.Option("missingkey=zero")
	}

	_, err := t.Parse(content)
	if err != nil {
		return err
	}

	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[key]; ok {
		return fmt.Errorf("template %s already exists", name)
	}

	r.m[key] = t
	return nil
}

func (r *Grender) getLayout(name string) *template.Template {
	r.RLock()
	defer r.RUnlock()
	return r.m["layout:"+name]
}

func (r *Grender) get(name string) *template.Template {
	r.RLock()
	defer r.RUnlock()
	return r.m["page:"+name]
}
