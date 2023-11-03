package rulestab

import (
	"sudonters/zootler/internal/entity"
	"sudonters/zootler/pkg/rules/parser"

	"github.com/charmbracelet/lipgloss"
)

type changeRuleView entity.Model
type addRuleView struct {
	id    entity.Model
	v     ruleView
	focus bool
}

type ruleView struct {
	name    string
	rawRule string
	ast     parser.Expression
	err     error
}

var ruleRowNameStyle = lipgloss.NewStyle().Width(15).Align(lipgloss.Left)

func (r ruleView) View() string {
	makeRow := func(name, content string) string {
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			ruleRowNameStyle.Render(name),
			"|",
			content,
		)
	}

	name := makeRow("Name", r.name)
	raw := makeRow("Raw Rule", r.rawRule)

	return lipgloss.JoinVertical(lipgloss.Left, name, "\n", raw)
}
