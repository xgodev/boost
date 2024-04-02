package cache

type Plugin[R any] interface {
	Del(*Context[R], string) error
	Get(*Context[R], string) ([]byte, error)
	Set(*Context[R], string, []byte) error
}
