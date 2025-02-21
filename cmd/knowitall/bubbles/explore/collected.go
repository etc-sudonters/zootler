package explore

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newCollected(sphere *NamedSphere) collected {
	var c collected
	var items []list.Item
	if sphere != nil {
		items = make([]list.Item, len(sphere.Tokens))
		for i := range sphere.Tokens {
			items[i] = listToken(sphere.Tokens[i])
			c.total += sphere.Tokens[i].Qty
		}
	}
	c.list = list.New(items, listTokenDelegate{}, 0, 0)
	listDefaults(&c.list)
	return c
}

type collected struct {
	total int
	list  list.Model
}

func (_ collected) Init() tea.Cmd { return nil }

func (this collected) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	add := func(c tea.Cmd) { cmds = append(cmds, c) }
	batch := func() tea.Cmd { return tea.Batch(cmds...) }
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.list.SetSize(msg.Width, msg.Height)
		return this, batch()
	}

	var cmd tea.Cmd
	this.list, cmd = this.list.Update(msg)
	add(cmd)

	return this, batch()
}

func (this collected) View() string {
	var view strings.Builder
	view.WriteString("INVENTORY\n")
	fmt.Fprintf(&view, "Collected: %3d\n", this.total)
	fmt.Fprint(&view, this.list.View())
	return view.String()
}

type listToken NamedToken

func (_ listToken) FilterValue() string {
	return ""
}

type listTokenDelegate struct{}

func (_ listTokenDelegate) Height() int  { return 1 }
func (_ listTokenDelegate) Spacing() int { return 1 }
func (_ listTokenDelegate) Update(tea.Msg, *list.Model) tea.Cmd {
	return nil
}
func (_ listTokenDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	token := item.(listToken)
	cursor := "  "
	if index == m.Index() {
		cursor = "> "
	}
	row := fmt.Sprintf("%3d %s", token.Qty, token.Name)
	fmt.Fprint(w, lipgloss.JoinHorizontal(lipgloss.Top, cursor, row))
}
