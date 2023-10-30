package foldertabs

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/etc-sudonters/substrate/bag"
)

type Opt func(*tabs)

func WithKeys(keys TabKeys) Opt {
	return func(t *tabs) {
		t.keys = keys
	}
}

func WithStyle(style TabStyle) Opt {
	return func(t *tabs) {
		t.style = style
	}
}

func WithTab(name string, m tea.Model) Opt {
	return func(t *tabs) {
		t.AddTab(name, m)
	}
}

func New(opts ...Opt) tabs {
	var t tabs
	t.names = make([]string, 4)
	t.content = make([]tea.Model, 4)
	t.keys = defaultTabKeys()

	for _, o := range opts {
		o(&t)
	}

	return t
}

type TabPosition int
type IsTabActive bool

const (
	TabPosStart TabPosition = iota
	TabPosMiddle
	TabPosEnd

	TabIsActive    IsTabActive = true
	TabIsNotActive IsTabActive = false
)

type TabStyle struct {
	Document, Window, Active, Inactive lipgloss.Style
}

type TabKeys struct {
	NextTab, PrevTab, CloseTab key.Binding
}

func defaultTabKeys() TabKeys {
	return TabKeys{
		NextTab:  key.NewBinding(key.WithKeys("ctrl+tab")),
		PrevTab:  key.NewBinding(key.WithKeys("ctrl+shift+tab")),
		CloseTab: key.NewBinding(key.WithKeys("ctrl+X")),
	}
}

type tabs struct {
	idx     int
	names   []string
	content []tea.Model

	keys  TabKeys
	style TabStyle
}

func (t *tabs) AddTab(name string, m tea.Model) {
	t.names = append(t.names, name)
	t.content = append(t.content, m)
}

func (t *tabs) CloseTab(name string) {
}

func (t *tabs) FocusOn(idx int) {
	if 0 <= idx && idx < len(t.names) {
		t.idx = idx
	}
}

func (t tabs) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, m := range t.content {
		cmds = append(cmds, m.Init())
	}
	return tea.Batch(cmds...)
}

func (t tabs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.keys.NextTab):
			t.idx = bag.Min(t.idx+1, len(t.names)-1)
			return t, nil
		case key.Matches(msg, t.keys.PrevTab):
			t.idx = bag.Max(t.idx-1, 0)
			return t, nil
		}
	}

	var cmd tea.Cmd
	var tab tea.Model

	tab = t.content[t.idx]
	tab, cmd = tab.Update(msg)
	t.content[t.idx] = tab

	return t, cmd
}

func (t tabs) View() string {
	var repr strings.Builder
	var rendered []string

	for i, name := range t.names {
		style := t.style.Inactive
		if i == t.idx {
			style = t.style.Active
		}

		rendered = append(rendered, style.Render(name))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	repr.WriteString(row)
	repr.WriteRune('\n')

	width := lipgloss.Width(row) - t.style.Window.GetHorizontalFrameSize()
	contentStyle := t.style.Window.Copy().Width(width)

	fmt.Fprint(&repr, contentStyle.Render(t.content[t.idx].View()))

	return t.style.Document.Render(repr.String())
}
