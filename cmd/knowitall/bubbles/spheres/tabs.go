package spheres

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func newSummary(parent *breakdown) summaryTab {
	return summaryTab{}
}

type summaryTab struct{}

func (_ *summaryTab) update(tea.Msg, *breakdown) tea.Cmd { return nil }
func (_ *summaryTab) init(*breakdown) tea.Cmd            { return nil }

func (_ summaryTab) view(content io.Writer, current *breakdown) {
	if current.sphere == nil {
		fmt.Fprintln(content, "no sphere loaded")
		return
	}
	fmt.Fprintf(content, "Breakdown of Sphere %03d\n", current.sphere.I)

	if current.sphere.Err != nil {
		fmt.Fprintf(content, "Sphere Error: %s\n", errStyle.Render(current.sphere.Err.Error()))
	}

	tmp := []struct {
		age   string
		nodes Handles
	}{
		{"Adult", current.sphere.Adult},
		{"Child", current.sphere.Child},
	}
	for _, pair := range tmp {
		age, nodes := pair.age, pair.nodes
		fmt.Fprintln(content, age)
		fmt.Fprintf(
			content,
			"\tTotal Visited: %03d\n\tPending Queue Size: %03d\n\tReached: %03d\n",
			nodes.Visited.Len(), nodes.Pending.Len(), nodes.Reached.Len(),
		)
	}

	fmt.Fprintln(content)
	fmt.Fprintf(content, "Total collected items: %03d", current.summary.Collected)

}

func newTokenTab(*breakdown) tokensTab {
	return tokensTab{}
}

type tokensTab struct{}

func (_ *tokensTab) update(tea.Msg, *breakdown) tea.Cmd { return nil }
func (_ *tokensTab) init(*breakdown) tea.Cmd            { return nil }

func (_ tokensTab) view(w io.Writer, current *breakdown) {
	if current.sphere == nil {
		fmt.Fprintln(w, "no sphere loaded")
		return
	}

	fmt.Fprintf(w, "Collected %3d tokens", len(current.sphere.Tokens))
}

func newSearchTab(*breakdown) runSearchTab {
	return runSearchTab{}
}

type runSearchTab struct{}

func (_ *runSearchTab) init(*breakdown) tea.Cmd { return nil }

func (_ *runSearchTab) update(msg tea.Msg, _ *breakdown) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return RequestSearch
		}
	}
	return nil
}

func (this runSearchTab) view(w io.Writer) {
	fmt.Fprint(w, "[Run Search]")
}

var breakDowntabs = []string{
	"SUMMARY",
	"TOKENS",
	"EDGES",
	"DISASSEMBLY",
	"SEARCH",
}

const (
	TAB_SUMMARY     = 0
	TAB_TOKENS      = 1
	TAB_EDGES       = 2
	TAB_DISASSEMBLY = 3
	TAB_SEARCH      = 4
)

func sizeTab(width int, tabCount int) int {
	return (width - tabCount) / tabCount
}

func renderBreakDownTabs(tabWidth int, current int) string {
	tabs := make([]string, len(breakDowntabs))
	for i, name := range breakDowntabs {
		first, last, active := i == 0, i == len(breakDowntabs)-1, i == current

		style := inactiveStyle.Border(inactiveTabBorder)
		if active {
			style = currentStyle.Border(activeTabBorder)
		}

		style = tabStyle.Inherit(style).Width(tabWidth).Height(1).AlignHorizontal(lipgloss.Center)

		border, _, _, _, _ := style.GetBorder()
		if first && active {
			border.BottomLeft = "│"
		} else if first && !active {
			border.BottomLeft = "├"
		} else if last && active {
			border.BottomRight = "│"
		} else if last && !active {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		tabs[i] = style.Render(name)
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	return row
}

var (
	tabStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	windowStyle       = lipgloss.NewStyle().Padding(2, 0).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	contentStyle      = lipgloss.NewStyle().Padding(0, 2)
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
