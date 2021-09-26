package grender

import "strings"

// Option is a callback to tune Grender.
type Option func(*Grender)

// OptionMissingKeyZero is to set templating engine with "missingkey=zero" option.
func OptionMissingKeyZero(yesno bool) Option {
	return func(r *Grender) {
		r.missingKeyZero = yesno
	}
}

// OptionTemplateExt sets extension name of template file for recognition. Default `.tmpl`.
func OptionTemplateExt(ext []string) Option {
	extMap := make(map[string]int)
	for _, e := range ext {
		if !strings.HasPrefix(e, ".") {
			e = "." + e
		}
		extMap[e] = 1
	}

	return func(r *Grender) {
		r.ext = extMap
	}
}

// OptionTemplateDir sets dir location of templates. Default `templates`.
func OptionTemplateDir(dir string) Option {
	return func(r *Grender) {
		r.templateDir = dir
	}
}

// OptionLayoutDir sets dir location of layout templates. Default `templates/layouts`.
func OptionLayoutDir(dir string) Option {
	return func(r *Grender) {
		r.layoutDir = dir
	}
}
