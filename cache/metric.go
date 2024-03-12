package cache

type Metric struct {
	Hit        int
	Miss       int
	GetSuccess int
	GetError   int
	DelSuccess int
	DelError   int
	SetSuccess int
	SetError   int
	Entries    int
	Driver     string
}
