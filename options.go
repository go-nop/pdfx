package pdfx

// Option is a function type for setting PDFProcessor options.
type Option func(*PDFProcessor)

// WithPassword is a function uses to set user password.
func WithPassword(password string) Option {
	return func(p *PDFProcessor) {
		configuration := p.configuration
		configuration.UserPW = password
		p.configuration = configuration
	}
}

// WithOptimize is a function uses to set optimize option.
func WithOptimize(optimize bool) Option {
	return func(p *PDFProcessor) {
		configuration := p.configuration
		configuration.Optimize = optimize
		p.configuration = configuration
	}
}
