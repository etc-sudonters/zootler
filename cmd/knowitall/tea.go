package main

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/etc-sudonters/substrate/dontio"
	"github.com/etc-sudonters/substrate/stageleft"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	std      *dontio.Std
}

func initialModel() model {
	return model{
		choices:  []string{"Do the hello world", "Take over the world"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		str := msg.String()
		m.std.WriteLineOut("processing %s", str)
		switch str {
		case "ctrl+c", "q":
			m.std.WriteLineErr("graceful shutdown requested")
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			m.std.WriteLineErr("cursor: %d", m.cursor)
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
			m.std.WriteLineErr("cursor: %d", m.cursor)
		case "enter":
			_, exists := m.selected[m.cursor]
			if exists {
				m.std.WriteLineErr("unchecking %d", m.cursor)
				delete(m.selected, m.cursor)
			} else {
				m.std.WriteLineErr("checking %d", m.cursor)
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Todo list\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, exists := m.selected[i]; exists {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"
	return s
}

func runMain(ctx context.Context, std *dontio.Std, _ *cliOptions) stageleft.ExitCode {
	initial := initialModel()
	initial.std = std
	p := tea.NewProgram(initial)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(std.Err, ":\\\n%s", err)
		return 1
	}
	return 0
}
