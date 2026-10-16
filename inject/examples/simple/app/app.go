// Package app consumes the values provided by the other packages.
package app

import (
	"sort"

	"github.com/xgodev/boost/inject/examples/simple/greeting"
)

// Board collects the lines written by Announce.
type Board struct {
	Lines []string
}

// NewBoard provides the board.
// @Provide(index=0)
func NewBoard() *Board {
	return &Board{}
}

// Announce greets every name of the group, formally and casually. The
// casual greeter and the board are resolved by type, without @Inject.
// @Inject(name=formal, index=0)
// @Inject(group=names, index=2)
// @Invoke
func Announce(formal *greeting.Greeter, casual *greeting.Greeter, names []string, board *Board) {
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)

	for _, name := range sorted {
		board.Lines = append(board.Lines, formal.Greet(name), casual.Greet(name))
	}
}
