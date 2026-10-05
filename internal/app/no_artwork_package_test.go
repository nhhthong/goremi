// Tests that the artwork package and the module that only it needed are gone from the repository.
package app

import (
	"os"
	"strings"
	"testing"
)

// START: TestNoArtworkPackage

func TestNoArtworkPackage(t *testing.T) {
	if _, err := os.Stat("../artwork"); !os.IsNotExist(err) {
		t.Fatalf("internal/artwork still exists (stat error: %v)", err)
	}
}

// END: TestNoArtworkPackage

// START: TestGoModHasNoImageModule

func TestGoModHasNoImageModule(t *testing.T) {
	mod, err := os.ReadFile("../../go.mod")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(mod), "golang.org/x/image") {
		t.Fatal("go.mod still requires golang.org/x/image")
	}
}

// END: TestGoModHasNoImageModule
