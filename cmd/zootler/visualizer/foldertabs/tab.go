package foldertabs

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/etc-sudonters/substrate/bag"
)

type BorderDelegate func(border *lipgloss.Border, active, first, last bool)

func AddTab(name string, m tea.Model, focused bool) tea.Cmd {
	return func() tea.Msg {
		return addTab{
			Name:  name,
			Model: m,
			Focus: focused,
		}
	}
}

func CloseTab(name string) tea.Cmd {
	return func() tea.Msg {
		return closeTab{
			Name: name,
		}
	}
}

type addTab struct {
	Name  string
	Model tea.Model
	Focus bool
}

type closeTab struct {
	Name string
}

type Opt func(*Model)

type Style struct {
	Active, Inactive, ModelView lipgloss.Style

	BorderDelegate BorderDelegate
	HorizontalFill string
}

type Keys struct {
	NextTab, PrevTab key.Binding
}

func WithTab(name string, m tea.Model) Opt {
	return func(t *Model) {
		t.content = append(t.content, m)
		t.names = append(t.names, name)
	}
}

func WithStyle(s Style) Opt {
	return func(t *Model) {
		t.style = s
	}
}

func WithKeys(k Keys) Opt {
	return func(t *Model) {
		t.keys = k
	}
}

func defaultKeys() Keys {
	return Keys{
		PrevTab: key.NewBinding(key.WithKeys("[")),
		NextTab: key.NewBinding(key.WithKeys("]")),
	}
}

func New(opts ...Opt) Model {
	var t Model
	t.keys = defaultKeys()

	for _, o := range opts {
		o(&t)
	}

	return t
}

type Model struct {
	names   []string
	content []tea.Model
	cur     int

	contentW, contentH int

	fullW, fullH int

	style Style
	keys  Keys
}

func (t Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range t.content {
		cmds = append(cmds, c.Init())
	}
	return tea.Batch(cmds...)
}

func (t Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		vertTabSpace := bag.Max(
			t.style.Active.GetVerticalFrameSize(),
			t.style.Inactive.GetVerticalFrameSize(),
		)
		t.fullW, t.fullH = msg.Width, msg.Height
		t.contentW = msg.Width - t.style.ModelView.GetHorizontalFrameSize()
		t.contentH = msg.Height - vertTabSpace - t.style.ModelView.GetVerticalFrameSize()
		t.style.ModelView = t.style.ModelView.Width(t.contentW).Height(t.contentH)

		cmds := make([]tea.Cmd, len(t.content))
		for i := range t.content {
			t.content[i], cmds[i] = t.content[i].Update(tea.WindowSizeMsg{
				Width:  t.contentW,
				Height: t.contentH,
			})
		}

		cmds = append(cmds, cmd)

		return t, tea.Batch(cmds...)
	case addTab:
		t.AddTab(msg.Name, msg.Model)
		if msg.Focus {
			t.FocusOn(msg.Name)
		}
		return t, cmd
	case closeTab:
		if len(t.names) == 1 {
			return t, tea.Quit
		}
		t.CloseTab(msg.Name)
		return t, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.keys.NextTab):
			t.cur = bag.Min(t.cur+1, len(t.names)-1)
			return t, cmd
		case key.Matches(msg, t.keys.PrevTab):
			t.cur = bag.Max(t.cur-1, 0)
			return t, cmd
		}
	}

	model := t.content[t.cur]
	model, cmd = model.Update(msg)
	t.content[t.cur] = model
	return t, cmd
}

func (t Model) View() string {
	var repr strings.Builder
	var rendered []string

	for i, tabName := range t.names {
		style := t.style.Inactive
		if i == t.cur {
			style = t.style.Active
		}

		active, first, last := t.cur == i, i == 0, i == len(t.names)-1
		border, _, _, _, _ := style.GetBorder()
		t.style.BorderDelegate(&border, active, first, last)

		style = style.Border(border)

		rendered = append(rendered, style.Render(tabName))
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	tabsWidth := lipgloss.Width(tabs)
	gap := strings.Repeat(t.style.HorizontalFill, bag.Max(0, t.fullW-tabsWidth-3))
	modelStyle := t.style.ModelView.Height(t.contentH - lipgloss.Height(tabs))
	modelBorder, _, _, _, _ := modelStyle.GetBorder()

	fmt.Fprint(&repr, lipgloss.JoinHorizontal(lipgloss.Bottom, tabs, gap, modelBorder.TopRight), "\n")
	fmt.Fprint(&repr, modelStyle.Render(t.content[t.cur].View()))

	return repr.String()
}

func (t *Model) SetWidth(w int) {
	if w < 1 {
		return
	}

	t.style.ModelView = t.style.ModelView.Width(w)
}

func (t *Model) AddTab(name string, m tea.Model) {
	t.names = append(t.names, name)
	t.content = append(t.content, m)
}

func (t *Model) CloseTab(name string) {
	idx := -1

	for i, n := range t.names {
		if n == name {
			idx = i
			break
		}
	}

	if idx > -1 {
		t.names = append(t.names[:idx], t.names[idx+1:]...)
		t.content = append(t.content[:idx], t.content[idx+1:]...)
		t.cur = bag.Max(bag.Min(t.cur, len(t.names)-1), 0)
	}
}

func (t *Model) FocusOn(name string) {
	for i, n := range t.names {
		if n == name {
			t.cur = i
			break
		}
	}
}
