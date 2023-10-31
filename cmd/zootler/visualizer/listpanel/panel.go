package listpanel

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectListEntry func(list.Item, int) tea.Cmd

type Opt func(m *Model)

func WithList(l list.Model) Opt {
	return func(m *Model) {
		m.l = l
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
	l     list.Model
	state panelstate

	rawH, rawW int

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
	l.l.Styles = l.styles.Active
}

func (l *Model) Blur() {
	l.state = panelblurred
	l.l.Styles = l.styles.Inactive
}

func (l *Model) SetSize(w, h int) {
	l.l.SetSize(w, h)
}

func (l *Model) SetHeight(h int) {
	l.l.SetHeight(h)
}

func (l *Model) SetWidth(w int) {
	l.l.SetWidth(w)
}

func (l *Model) Height() int {
	return l.l.Height()
}

func (l *Model) Width() int {
	return l.l.Width()
}

func (l Model) Init() tea.Cmd {
	return nil
}

func (l Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if l.state == panelblurred {
		return l, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.l.SetSize(msg.Width, msg.Height)
		return l, nil
	}

	var cmd tea.Cmd
	l.l, cmd = l.l.Update(msg)

	return l, cmd
}

func (l Model) View() string {
	return l.l.View()
}

func (l Model) Title() string {
	return l.l.Title
}
