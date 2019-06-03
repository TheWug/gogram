package gogram

import (
	"strings"
)

// A helper function which replaces em-dashes with double dashes.
// this reverses a substition some telegram clients perform automatically
// and should be used when performing argparse-like command processing.
func NormalizeDashes(cmd string) (string) {
	return strings.Replace(cmd, "—", "--", -1)
}
