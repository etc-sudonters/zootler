package application

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Opt func(*application)

type GlobalKey struct {
	key.Binding
	Cmd tea.Cmd
}

func WithGlobalKeys(keys []GlobalKey) Opt {
	return func(a *application) {
		a.keys = append(a.keys, keys...)
	}
}

func WithStyle(s lipgloss.Style) Opt {
	return func(a *application) {
		a.style = s
	}
}

func New(m tea.Model, opts ...Opt) application {
	var app application
	app.m = m
	app.keys = append(app.keys, GlobalKey{
		Binding: key.NewBinding(key.WithKeys("ctrl+c")),
		Cmd:     tea.Quit,
	})
	for _, o := range opts {
		o(&app)
	}

	return app
}

type application struct {
	m     tea.Model
	keys  []GlobalKey
	style lipgloss.Style
}

func (a application) Init() tea.Cmd {
	return a.m.Init()
}

func (a application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.m, cmd = a.m.Update(tea.WindowSizeMsg{
			Height: msg.Height - a.style.GetVerticalFrameSize(),
			Width:  msg.Width - a.style.GetHorizontalFrameSize(),
		})
		return a, cmd
	case tea.KeyMsg:
		for _, k := range a.keys {
			if key.Matches(msg, k.Binding) {
				return a, k.Cmd
			}
		}
	}

	a.m, cmd = a.m.Update(msg)
	return a, cmd
}

func (a application) View() string {
	return a.style.Render(a.m.View())
}
