package tui

import (
	"sudonters/zootler/internal/tui/foldertabs"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var tabKeys = foldertabs.Keys{
	NextTab:  key.NewBinding(key.WithKeys("["), key.WithHelp("[", "Next Tab")),
	PrevTab:  key.NewBinding(key.WithKeys("]"), key.WithHelp("]", "Previous Tab")),
	CloseTab: key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "Close Tab")),
}

var tabStyle = foldertabs.Style{
	Active:    lipgloss.NewStyle().Padding(0, 0).Border(tabBorderWithBottom("┘", " ", "└"), true),
	Inactive:  lipgloss.NewStyle().Padding(0, 0).Border(tabBorderWithBottom("┴", "─", "┴"), true),
	ModelView: lipgloss.NewStyle().Padding(1).Border(lipgloss.NormalBorder()).UnsetBorderTop(),
	BorderDelegate: func(border *lipgloss.Border, active, first, last bool) {
		if first && active {
			border.BottomLeft = "│"
		} else if first && !active {
			border.BottomLeft = "├"
		} else if last && active {
			border.BottomRight = "│"
		} else if last && !active {
			border.BottomRight = "┤"
		}
	},
	HorizontalFill: "─",
}
