// package grender is a smart templating tool.
package grender

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

// Grender is the object to handle the rendering work.
type Grender struct {
	lock           sync.RWMutex
	m              map[string]*template.Template
	ext            map[string]int
	missingKeyZero bool
	templateDir    string
	layoutDir      string
}

// New creates one Grender object.
func New(opts ...Option) *Grender {
	r := &Grender{
		m:           make(map[string]*template.Template),
		ext:         map[string]int{".tmpl": 1},
		templateDir: "templates",
		layoutDir:   "templates/layouts",
	}

	for _, opt := range opts {
		opt(r)
	}
	return r
}

// LoadFromFS loads templates (include layouts) from file system.
func (r *Grender) LoadFromFS(fsys fs.FS) error {
	if r.layoutDir != "" {
		if err := r.load(fsys, true); err != nil {
			return err
		}
	}

	if r.templateDir != "" {
		if err := r.load(fsys, false); err != nil {
			return err
		}
	}

	return nil
}

func (r *Grender) load(fsys fs.FS, isLayout bool) error {
	var dir string
	if isLayout {
		dir = r.layoutDir
	} else {
		dir = r.templateDir
	}

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		ext := filepath.Ext(info.Name())
		if _, ok := r.ext[ext]; !ok {
			continue
		}

		content, err := fs.ReadFile(fsys, filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(filepath.Base(info.Name()), ext)
		if err := r.add(name, string(content), isLayout); err != nil {
			return err
		}
	}

	return nil
}

// AddLayout adds a layout template.
func (r *Grender) AddLayout(name, content string) error {
	return r.add(name, content, true)
}

// Add adds a normal (page) template.
func (r *Grender) Add(name, content string) error {
	return r.add(name, content, false)
}

func (r *Grender) add(name, content string, isLayout bool) error {
	var key string
	if isLayout {
		key = "layout:" + name
	} else {
		key = "page:" + name
	}

	t := template.New(key).Funcs(template.FuncMap{
		"current": func() string {
			return name
		},
		"share": func(sharedTempleteName string) string {
			return fmt.Sprintf("{{ share \"%s\" }}", sharedTempleteName)
		},
	})

	if isLayout {
		t.Funcs(template.FuncMap{
			"yield": func() string {
				return "{{.}}"
			},
		})
	}

	if r.missingKeyZero {
		t.Option("missingkey=zero")
	}

	_, err := t.Parse(content)
	if err != nil {
		return err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.m[key]; ok {
		return fmt.Errorf("template %s already exists", name)
	}

	r.m[key] = t
	return nil
}

func (r *Grender) get(name string, isLayout bool) *template.Template {
	var key string
	if isLayout {
		key = "layout:" + name
	} else {
		key = "page:" + name
	}

	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.m[key]
}
