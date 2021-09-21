// package grender is a smart templating tool.
package grender

import (
	"fmt"
	"sync"
	"text/template"
)

type Grender struct {
	sync.RWMutex
	m map[string]*template.Template
}

func NewGrender() *Grender {
	return &Grender{
		m: make(map[string]*template.Template),
	}
}

func (Grender *Grender) AddLayout(name, content string) error {
	key := "layout:" + name
	t, err := template.New(key).Option("missingkey=zero").Funcs(template.FuncMap{
		"yield": func() string {
			return "{{.}}"
		},
		"current": func() string {
			return name
		},
		"share": func(sharedTempleteName string) string {
			return fmt.Sprintf("{{ share \"%s\" }}", sharedTempleteName)
		},
	}).Parse(content)

	if err != nil {
		return err
	}

	Grender.Lock()
	defer Grender.Unlock()

	if _, ok := Grender.m[key]; ok {
		return fmt.Errorf("template %s already exists", name)
	}

	Grender.m[key] = t
	return nil
}

func (Grender *Grender) Add(name, content string) error {
	key := "page:" + name
	t, err := template.New(key).Option("missingkey=zero").Funcs(template.FuncMap{
		"current": func() string {
			return name
		},
		"share": func(sharedTempleteName string) string {
			return fmt.Sprintf("{{ share \"%s\" }}", sharedTempleteName)
		},
	}).Parse(content)

	if err != nil {
		return err
	}

	Grender.Lock()
	defer Grender.Unlock()

	if _, ok := Grender.m[key]; ok {
		return fmt.Errorf("template %s already exists", name)
	}

	Grender.m[key] = t
	return nil
}

func (Grender *Grender) GetLayout(name string) *template.Template {
	Grender.RLock()
	defer Grender.RUnlock()
	return Grender.m["layout:"+name]
}

func (Grender *Grender) Get(name string) *template.Template {
	Grender.RLock()
	defer Grender.RUnlock()
	return Grender.m["page:"+name]
}
