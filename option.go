package grender

// Option is a callback to tune Grender.
type Option func(*Grender)

// OptionMissingKeyZero is to set templating engine with "missingkey=zero" option.
func OptionMissingKeyZero(yesno bool) Option {
	return func(r *Grender) {
		r.missingKeyZero = yesno
	}
}
