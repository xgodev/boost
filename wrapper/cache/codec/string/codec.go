package string

type Codec[T string] struct {
}

func (c *Codec[T]) Encode(data T) (b []byte, err error) {
	return []byte(data), err
}

func (c *Codec[T]) Decode(b []byte, data *string) (err error) {
	*data = string(b)
	return nil
}

func New[T string]() *Codec[T] {
	return &Codec[T]{}
}
