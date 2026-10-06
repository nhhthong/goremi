// Tests of .github/workflows/release.yml: when it runs, what it may do, and what it publishes.
package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// START: TestReleaseWorkflowShape

func TestReleaseWorkflowShape(t *testing.T) {
	text, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(text)
	for _, want := range []string{"tags: ['v*']", "contents: write", "sh scripts/release-build.sh", "gh release create", "dist/*", "GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}"} {
		if !strings.Contains(s, want) {
			t.Errorf("release.yml lacks %q", want)
		}
	}
	if strings.Contains(s, "branches") || strings.Contains(s, "pull_request") {
		t.Error("release.yml must run on a pushed tag only, no branches or pull_request")
	}
}

// END: TestReleaseWorkflowShape
