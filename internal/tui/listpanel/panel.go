package listpanel

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectListEntry func(list.Item, int) tea.Cmd

type Opt func(m *Model)

func WithList(l list.Model) Opt {
	return func(m *Model) {
		m.L = l
	}
}

func WithStyles(s Styles) Opt {
	return func(m *Model) {
		m.styles = s
	}
}

type Styles struct {
	Active   list.Styles
	Inactive list.Styles
}

type panelstate int

const (
	panelblurred panelstate = iota
	panelfocused
)

type Model struct {
	L     list.Model
	state panelstate

	styles Styles
}

func New(opts ...Opt) Model {
	var m Model
	for _, o := range opts {
		o(&m)
	}

	return m
}

func (l *Model) Focus() {
	l.state = panelfocused
	l.L.Styles = l.styles.Active
}

func (l *Model) Blur() {
	l.state = panelblurred
	l.L.Styles = l.styles.Inactive
}

func (l Model) Init() tea.Cmd {
	return nil
}

func (l Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		l.L.SetSize(size.Width, size.Height)
		return l, nil
	}

	if l.state == panelblurred {
		return l, nil
	}

	var cmd tea.Cmd
	l.L, cmd = l.L.Update(msg)

	return l, cmd
}

func (l Model) View() string {
	return l.L.View()
}
