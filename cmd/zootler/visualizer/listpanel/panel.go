package listpanel

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectListEntry func(list.Item, int) tea.Cmd

type panelstate int

const (
	panelblurred panelstate = iota
	panelfocused
)

type Model struct {
	l     list.Model
	state panelstate
}

func New(
	l list.Model,
) Model {
	return Model{
		l: l,
	}
}

func (l *Model) Focus() {
	l.state = panelfocused
}

func (l *Model) Blur() {
	l.state = panelblurred
}

func (l Model) Init() tea.Cmd {
	return nil
}

func (l Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if l.state == panelblurred {
		return l, nil
	}
	var cmd tea.Cmd
	l.l, cmd = l.l.Update(msg)

	return l, cmd
}

func (l Model) View() string {
	return l.l.View()
}
