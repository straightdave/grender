package grender

import (
	"bytes"
	"fmt"
	"text/template"
)

func (r *Grender) Render(layoutName, pageName string, data interface{}) (string, error) {
	if layoutName == "" {
		return r.renderWithoutLayout(pageName, data)
	}
	return r.renderWithLayout(layoutName, pageName, data)
}

func (r *Grender) renderWithLayout(layoutName, pageName string, data interface{}) (string, error) {
	var b1 bytes.Buffer
	layout := r.GetLayout(layoutName)
	if layout == nil {
		return "", fmt.Errorf("no template %s", layoutName)
	}
	if err := layout.Execute(&b1, data); err != nil {
		return "", err
	}
	layoutContent := string(b1.Bytes())

	pageContent, err := r.renderWithoutLayout(pageName, data)
	if err != nil {
		return "", err
	}

	outTmpl, err := template.New("output-template").Funcs(template.FuncMap{
		"share": func(sharedTempleteName string) string {
			shared := r.Get(sharedTempleteName)
			if shared == nil {
				panic(fmt.Errorf("no shared template %s", sharedTempleteName))
			}

			var b bytes.Buffer
			if err := shared.Execute(&b, data); err != nil {
				panic(err)
			}

			return string(b.Bytes())
		},
	}).Parse(layoutContent)
	if err != nil {
		return "", err
	}

	var b3 bytes.Buffer
	if err := outTmpl.Execute(&b3, pageContent); err != nil {
		return "", err
	}

	return string(b3.Bytes()), nil
}

func (r *Grender) renderWithoutLayout(pageName string, data interface{}) (string, error) {
	var b bytes.Buffer
	pt := r.Get(pageName)
	if pt == nil {
		return "", fmt.Errorf("no template %s", pageName)
	}
	if err := pt.Execute(&b, data); err != nil {
		return "", err
	}
	content := string(b.Bytes())

	// for shared Grender
	// if we know how to detect any "share" function, we can ignore this step

	temp, err := template.New("temp-template").Funcs(template.FuncMap{
		"share": func(sharedTempleteName string) string {
			shared := r.Get(sharedTempleteName)
			if shared == nil {
				panic(fmt.Errorf("no shared template %s", sharedTempleteName))
			}

			var b bytes.Buffer
			if err := shared.Execute(&b, data); err != nil {
				panic(err)
			}

			return string(b.Bytes())
		},
	}).Parse(content)
	if err != nil {
		return "", err
	}

	var bt bytes.Buffer
	if err := temp.Execute(&bt, data); err != nil {
		return "", err
	}
	return string(bt.Bytes()), nil
}
