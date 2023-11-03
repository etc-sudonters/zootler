package application

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func WantWindowSize() tea.Msg {
	return wantWindowSize{}
}

type wantWindowSize struct{}

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

func WithKeyDisplay() Opt {
	return func(a *application) {
		a.showKeys = true
	}
}

func WithMsgDisplay() Opt {
	return func(a *application) {
		a.showMsgs = true
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
	m        tea.Model
	size     tea.WindowSizeMsg
	keys     []GlobalKey
	style    lipgloss.Style
	showKeys bool
	lastKey  tea.KeyMsg
	lastMsg  tea.Msg
	showMsgs bool
}

func (a application) Init() tea.Cmd {
	return a.m.Init()
}

func (a application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if a.showMsgs {
		a.lastMsg = msg
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case wantWindowSize:
		if a.size.Height != 0 && a.size.Width != 0 {
			return a, func() tea.Msg { return a.size }
		}
	case tea.WindowSizeMsg:
		a.size = msg
		a.m, cmd = a.m.Update(tea.WindowSizeMsg{
			Height: msg.Height - a.style.GetVerticalFrameSize() - 1,
			Width:  msg.Width - a.style.GetHorizontalFrameSize(),
		})
		return a, cmd
	case tea.KeyMsg:
		if a.showKeys {
			a.lastKey = msg
		}
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
	var debugPieces []string

	content := a.style.Render(a.m.View())

	debugStyle := lipgloss.NewStyle().Width(lipgloss.Width(content) / 4).Align(lipgloss.Left)

	if a.showKeys {

		debugPieces = append(debugPieces, debugStyle.Render(fmt.Sprintf("Key:\t%s", a.lastKey.String())))
	}

	if a.showMsgs {
		debugPieces = append(debugPieces, debugStyle.Render(fmt.Sprintf("Msg:\t%T%v", a.lastMsg, a.lastMsg)))
	}

	debugBar := lipgloss.JoinHorizontal(lipgloss.Top, debugPieces...)

	return lipgloss.JoinVertical(lipgloss.Left, content, debugBar)
}

func ResizeWindowMsg(size tea.WindowSizeMsg, width, height float64) tea.WindowSizeMsg {
	return tea.WindowSizeMsg{
		Width:  int(float64(size.Width) * width),
		Height: int(float64(size.Height) * height),
	}
}
