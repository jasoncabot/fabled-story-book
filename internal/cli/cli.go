package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jasoncabot/fabled-story-book/internal/jabl"
	"github.com/muesli/termenv"
)

type jablProgram struct {
	interpreter      *jabl.Interpreter
	currentSection   jabl.SectionId
	currentSelection int
	profile          termenv.Profile
	text             string
	choices          []*jabl.Choice
	transition       jabl.SectionId
	err              error
}

func NewJABLProgram(i *jabl.Interpreter, entrypoint string, profile termenv.Profile) *jablProgram {
	return &jablProgram{
		interpreter:      i,
		currentSection:   jabl.SectionId(entrypoint),
		currentSelection: 0,
		profile:          profile,
		text:             "",
		choices:          []*jabl.Choice{},
		transition:       "",
	}
}

func (p *jablProgram) Init() tea.Cmd {
	// When starting we want to execute the current section
	return execSection(p.interpreter, p.currentSection)
}

func (p *jablProgram) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Check for keypress
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return p, tea.Quit
		case "up", "k":
			if p.currentSelection > 0 {
				p.currentSelection--
			}
		case "down", "j":
			if p.currentSelection < len(p.choices)-1 {
				p.currentSelection++
			}
		case "enter", " ":
			// Find the section id for the selected choice
			if p.currentSelection < len(p.choices) {
				choice := p.choices[p.currentSelection]
				return p, execChoice(p.interpreter, choice)
			}
			return p, nil
		}

	case executeMsg:
		if msg.result != nil {
			p.choices = msg.result.Choices
			p.text = msg.result.Output
			p.transition = msg.result.Transition
		} else {
			p.choices = []*jabl.Choice{}
			p.text = ""
			p.transition = ""
		}
		p.err = msg.err
		return p, nil

	}

	return p, nil
}

func (p *jablProgram) View() string {
	// The header
	s := p.text + "\n\n"

	if p.err == nil {
		// Iterate over our choices
		choices := p.choices
		for i, choice := range choices {
			if p.currentSelection == i {
				s += fmt.Sprintf("  %s\n", p.profile.String(choice.Text).Underline().String())
			} else {
				s += fmt.Sprintf("  %s\n", choice.Text)
			}
		}
	} else {
		s += p.profile.String(fmt.Sprintf("Error: %v\n", p.err)).Foreground(termenv.ANSIBrightRed).String()
	}

	// The footer
	s += p.profile.String("\nPress q to quit.\n").Faint().String()

	// Send the UI for rendering
	return s
}
