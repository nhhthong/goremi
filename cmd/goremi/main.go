// Entry point of the goremi binary: hands the arguments to app.Run with the YouTube provider.
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
	path, _ := app.ConfigPath() // no config dir leaves path empty, which means the defaults
	run := func(m tea.Model) (tea.Model, error) { return tea.NewProgram(m).Run() }
	if err := app.Run(os.Args[1:], path, &provider.YouTubeProvider{}, run); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// END: main
