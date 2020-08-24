package Options

var DefaultCamelCase = true

type LoadEnvOptions struct {
	CamelCase *bool
}

func LoadEnv() *LoadEnvOptions {
	return &LoadEnvOptions{
		CamelCase: &DefaultCamelCase,
	}
}

func (leo *LoadEnvOptions) SetCamelCase(b bool) *LoadEnvOptions {
	leo.CamelCase = &b
	return leo
}

// MergeLoadEnvOptions combines the given *LoadEnvOptions into one *LoadEnvOptions in a last one wins fashion.
func MergeLoadEnvOptions(opts ...*LoadEnvOptions) *LoadEnvOptions {
	uOpts := LoadEnv()
	for _, uo := range opts {
		if uo == nil {
			continue
		}
		if uo.CamelCase != nil {
			uOpts.CamelCase = uo.CamelCase
		}
	}

	return uOpts
}
