// Helpers shared by the installer tests of every system.
package install

import (
	"os"
	"strings"
)

// START: shared

// logLines returns the non-empty lines of a file, none when it does not exist.
func logLines(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var out []string
	for _, l := range strings.Split(string(data), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// has tells whether some line of lines equals want.
func has(lines []string, want string) bool {
	for _, l := range lines {
		if l == want {
			return true
		}
	}
	return false
}

// END: shared
