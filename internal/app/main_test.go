// Test setup of the app package: the tests run with a private cache directory, so Run never writes the user's log.
package app

import (
	"os"
	"testing"
)

// START: TestMain

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "goremi-test-cache-")
	if err != nil {
		panic(err)
	}
	os.Setenv("XDG_CACHE_HOME", dir)
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// END: TestMain
