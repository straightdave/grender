package grender

import (
	"bytes"
	"fmt"
	"text/template"
)

// Result is the result of rendering, actually a byte slice.
type Result []byte

// String converts Result into a byte slice.
func (r Result) Bytes() []byte {
	return r
}

// String converts Result into a String.
func (r Result) String() string {
	return string(r)
}

// Render renders with given layout and page template names, as well as data.
func (r *Grender) Render(layoutName, pageName string, data interface{}) (Result, error) {
	if layoutName == "" {
		return r.renderWithoutLayout(pageName, data)
	}
	return r.renderWithLayout(layoutName, pageName, data)
}

func (r *Grender) renderWithLayout(layoutName, pageName string, data interface{}) (Result, error) {
	var b1 bytes.Buffer
	layout := r.get(layoutName, true)
	if layout == nil {
		return b1.Bytes(), fmt.Errorf("no template %s", layoutName)
	}
	if err := layout.Execute(&b1, data); err != nil {
		return b1.Bytes(), err
	}
	layoutContent := b1.String()

	pageContent, err := r.renderWithoutLayout(pageName, data)
	if err != nil {
		return b1.Bytes(), err
	}

	outTmpl, err := template.New("output-template").Funcs(template.FuncMap{
		"share": func(sharedTempleteName string) string {
			return r.renderShared(sharedTempleteName, data).String()
		},
	}).Parse(layoutContent)
	if err != nil {
		return b1.Bytes(), err
	}

	var b3 bytes.Buffer
	if err := outTmpl.Execute(&b3, pageContent); err != nil {
		return b3.Bytes(), err
	}

	return b3.Bytes(), nil
}

func (r *Grender) renderWithoutLayout(pageName string, data interface{}) (Result, error) {
	var b bytes.Buffer
	pt := r.get(pageName, false)
	if pt == nil {
		return b.Bytes(), fmt.Errorf("no template %s", pageName)
	}
	if err := pt.Execute(&b, data); err != nil {
		return b.Bytes(), err
	}
	content := b.String()

	// for shared Grender
	// if we know how to detect any "share" function, we can ignore this step

	temp, err := template.New("temp-template").Funcs(template.FuncMap{
		"share": func(sharedTempleteName string) string {
			return r.renderShared(sharedTempleteName, data).String()
		},
	}).Parse(content)
	if err != nil {
		return Result(b.Bytes()), err
	}

	var bt bytes.Buffer
	if err := temp.Execute(&bt, data); err != nil {
		return bt.Bytes(), err
	}
	return Result(bt.Bytes()), nil
}

func (r *Grender) renderShared(sharedTemplateName string, data interface{}) Result {
	shared := r.get(sharedTemplateName, false)
	if shared == nil {
		panic(fmt.Errorf("no shared template %s", sharedTemplateName))
	}

	var b bytes.Buffer
	if err := shared.Execute(&b, data); err != nil {
		panic(err)
	}
	return Result(b.Bytes())
}
