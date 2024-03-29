package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jasoncabot/fabled-story-book/internal/cli"
	"github.com/jasoncabot/fabled-story-book/internal/jabl"
	"github.com/muesli/termenv"
)

func main() {

	// load a particular character
	loader := cli.NewFileLoader(filepath.Join("assets", "example01"))

	interpreter := jabl.NewInterpreter(loader)

	term := termenv.EnvColorProfile()

	state := cli.NewStateMapper()
	session := cli.NewJABLProgram(interpreter, state, "entrypoint.jabl", term)

	app := tea.NewProgram(session)
	if _, err := app.Run(); err != nil {
		fmt.Printf("An error occurred: %v", err)
		os.Exit(1)
	}
}
