// Tests of the badge and the header of the empty screen: the mascot with the badge beside it or under it (ui tasks 3.4.3, 3.4.6, 3.4.7).
package ui

import (
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

// plainRows is the header without colour codes, split in rows.
func plainRows(s string) []string { return strings.Split(ansiCode.ReplaceAllString(s, ""), "\n") }

// col is the column of the first rune of want in row, or -1.
func col(row, want string) int {
	i := strings.Index(row, want)
	if i < 0 {
		return -1
	}
	return utf8.RuneCountInString(row[:i])
}

// START: TestBadgeLines

func TestBadgeLines(t *testing.T) {
	want := []string{"Goremi v0.1.1", "Created by nhhthong", "Current provider: YouTube"}
	if got := Badge("0.1.1"); !reflect.DeepEqual(got, want) {
		t.Fatalf("Badge(0.1.1) = %q, want %q", got, want)
	}
}

// END: TestBadgeLines

// START: TestBadgeUsesVersion

func TestBadgeUsesVersion(t *testing.T) {
	if got := Badge("dev")[0]; got != "Goremi vdev" {
		t.Fatalf("first badge line = %q, want Goremi vdev", got)
	}
}

// END: TestBadgeUsesVersion
