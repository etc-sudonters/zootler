package visualizer

import (
	"fmt"
	"strings"
	"sudonters/zootler/internal/entity/componenttable"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/etc-sudonters/substrate/bag"
)

type menuKeys struct {
	up, down, choose key.Binding
}

type menu struct {
	keys       menuKeys
	components []componenttable.RowData
	idx        int
}

func (m menu) Init() tea.Cmd {
	return nil
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.up):
			m.idx = bag.Max(1+m.idx, len(m.components)-1)
			break
		case key.Matches(msg, m.keys.down):
			m.idx = bag.Min(m.idx-1, 0)
			break
		case key.Matches(msg, m.keys.choose):
			break
		}
		break
	}

	var cmd tea.Cmd
	switch len(cmds) {
	case 0:
		break
	case 1:
		cmd = cmds[0]
		break
	default:
		cmd = tea.Batch(cmds...)
	}

	return m, cmd
}

func (m menu) View() string {
	var repr strings.Builder

	for idx, r := range m.components {
		if idx == m.idx {
			(&repr).WriteString(" x ")
		} else {
			(&repr).WriteString("   ")
		}
		fmt.Fprintf(&repr, "%s (%d/%d)\n", r.Type(), r.Len(), r.Capacity())
	}

	return repr.String()
}
