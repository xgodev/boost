package cache

type OptionSet func() func(opt *Option)

func SaveEmpty() func(opt *Option) {
	return func(opt *Option) {
		opt.SaveEmpty = true
	}
}

func AsyncSave() func(opt *Option) {
	return func(opt *Option) {
		opt.AsyncSave = true
	}
}

func WithoutReplicate() func(opt *Option) {
	return func(opt *Option) {
		opt.Replicate = false
	}
}

type Option struct {
	SaveEmpty bool
	AsyncSave bool
	Replicate bool
}
