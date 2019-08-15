// lexresolver.go

package crdt

import (
	"fmt"
	"strings"

	"github.com/nsip/vvmap"
)

//
// conflict resolver to choose between versions of a value when
// versions are the same
//
// needs Go 12, where maps are printed with all elements in key order.
//
var LexicographicConflictResolver = func(key string, left, right vvmap.Record) bool {
	leftVal := fmt.Sprintf("%v", left.Value)
	rightVal := fmt.Sprintf("%v", right.Value)
	return strings.Compare(leftVal, rightVal) > 0 // choose left if lexicographically greater
}

var DefaultResolver = LexicographicConflictResolver
