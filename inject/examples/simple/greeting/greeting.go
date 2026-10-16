// Package greeting provides the greeters used by the example.
package greeting

// Greeter builds greeting lines.
type Greeter struct {
	Prefix string
}

// Greet returns the greeting for name.
func (g *Greeter) Greet(name string) string {
	return g.Prefix + ", " + name + "!"
}

// NewGreeter provides the default greeter.
// @Provide(index=0)
func NewGreeter() *Greeter {
	return &Greeter{Prefix: "Hello"}
}

// NewFormalGreeter provides the greeter named "formal".
// @Provide(name=formal, index=0)
func NewFormalGreeter() (*Greeter, error) {
	return &Greeter{Prefix: "Good morning"}, nil
}
