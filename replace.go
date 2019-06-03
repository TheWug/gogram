package gogram

import (
	"strings"
)

// A helper function which replaces leading em-dashes with double dashes.
// this reverses a substition some telegram clients perform automatically
// and should be used when performing argparse-like command processing.
func NormalizeDashes(cmd string) (string) {
	new_cmd := strings.TrimPrefix(cmd, "—") // remove em dash
	if cmd != new_cmd {
		// if the original string WAS prefixed with an em dash, add a double dash instead.
		return "--" + new_cmd
	} else {
		// otherwise, return the original string unchanged.
		return cmd
	}
}
