package binary

// Define a constraint that T must be []byte
type ByteArray interface {
	~[]byte // Ensures only []byte is accepted
}

type Codec[T ByteArray] struct {
}

func (c *Codec[T]) Encode(e T) (b []byte, err error) {
	// Since T is already []byte, we can return it directly without any further encoding
	return e, nil
}

func (c *Codec[T]) Decode(b []byte, data *T) (err error) {
	// Decode is essentially just assigning the byte array
	*data = b
	return nil
}

func New[T ByteArray]() *Codec[T] {
	return &Codec[T]{}
}
