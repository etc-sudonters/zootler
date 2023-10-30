package visualizer

import (
	"context"
	"sudonters/zootler/cmd/zootler/visualizer/application"
	"sudonters/zootler/cmd/zootler/visualizer/foldertabs"
	"sudonters/zootler/cmd/zootler/visualizer/panels"
	"sudonters/zootler/internal/entity/bitpool"
	"sudonters/zootler/internal/reitertools"
	"sudonters/zootler/pkg/world"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Visualize(w world.World) visualizer {
	return visualizer{w}
}

type visualizer struct {
	w world.World
}

func (v visualizer) Run(ctx context.Context) error {
	tbl, err := bitpool.ExtractComponentTable(v.w.Entities.Pool)
	if err != nil {
		return err
	}

	rows := reitertools.ToSlice(tbl.Rows())
	componentPanel := createComponentOverview(rows, v.w.Entities.Pool)
	(&componentPanel).Focus()

	p := panels.New(componentPanel)

	tabHighlight := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	tabs := foldertabs.New(
		foldertabs.WithTab("Components", p),
		foldertabs.WithTab("Other stuff", other("stuff")),
		foldertabs.WithStyle(foldertabs.TabStyle{
			Inactive: lipgloss.NewStyle().Border(tabBorderWithBottom("┴", "─", "┴"), true).
				BorderForeground(tabHighlight).Padding(0, 1),

			Active: lipgloss.NewStyle().Border(tabBorderWithBottom("┘", " ", "└"), true).
				BorderForeground(tabHighlight).Padding(0, 1),

			Document: lipgloss.NewStyle().Padding(1, 2, 1, 2),

			Window: lipgloss.NewStyle().BorderForeground(tabHighlight).
				Padding(2, 0).
				Align(lipgloss.Center).
				Border(lipgloss.NormalBorder()).
				UnsetBorderTop(),
		}),
	)
	tabs.FocusOn(1)

	app := application.New(p)
	_, err = tea.NewProgram(app, tea.WithAltScreen(), tea.WithContext(ctx)).Run()
	return err
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

type other string

func (o other) Init() tea.Cmd {
	return nil
}

func (o other) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return o, nil
}

func (o other) View() string {
	return string(o)
}
