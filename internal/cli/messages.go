package cli

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jasoncabot/fabled-story-book/internal/jabl"
)

type executeMsg struct {
	result *jabl.Result
	err    error
}

func execSection(interpreter *jabl.Interpreter, id jabl.SectionId) func() tea.Msg {
	return func() tea.Msg {
		var msg tea.Msg

		wg := sync.WaitGroup{}
		wg.Add(1)
		interpreter.Execute(id, func(res *jabl.Result, err error) {
			msg = executeMsg{
				result: res,
				err:    err,
			}
			wg.Done()
		})

		wg.Wait()

		return msg
	}
}

func execChoice(interpreter *jabl.Interpreter, choice *jabl.Choice) func() tea.Msg {
	return func() tea.Msg {
		var msg executeMsg

		wg := sync.WaitGroup{}
		wg.Add(1)
		interpreter.Evaluate(choice.Code, func(res *jabl.Result, err error) {
			msg = executeMsg{
				result: res,
				err:    err,
			}
			wg.Done()
		})

		wg.Wait()

		return msg
	}
}
