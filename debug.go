package github

import (
	"fmt"
	"os"
)

func debugln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}
