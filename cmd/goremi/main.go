// Entry point of the goremi binary: runs the application with the YouTube provider.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"goremi/internal/app"
	"goremi/internal/provider"
)

// START: main

func main() {
	p := tea.NewProgram(app.New(&provider.YouTubeProvider{}))
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// END: main
