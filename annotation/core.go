package annotation

import "strings"

var (
	corePkgs = []string{
		"golang.org",
		"archive",
		"zip",
		"bufio",
		"builtin",
		"bytes",
		"compress",
		"container",
		"context",
		"crypto",
		"database",
		"debug",
		"embed",
		"encoding",
		"errors",
		"expvar",
		"flag",
		"fmt",
		"go",
		"hash",
		"html",
		"image",
		"index",
		"io",
		"log",
		"math",
		"mime",
		"net",
		"os",
		"path",
		"plugin",
		"reflect",
		"regexp",
		"runtime",
		"sort",
		"strconv",
		"strings",
		"sync",
		"syscall",
		"testing",
		"text",
		"time",
		"unicode",
		"unsafe",
	}
)

func isCorePackage(pkgPath string) bool {
	for _, n := range corePkgs {
		if strings.HasPrefix(pkgPath, n) {
			return true
		}
	}
	return false
}
